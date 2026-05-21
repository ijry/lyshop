package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ijry/lyshop/core/db"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"gorm.io/gorm"
)

// CreateOrderReq is the request payload to create an order.
type CreateOrderReq struct {
	UserID        uint64          `json:"user_id"`
	AddressID     uint64          `json:"address_id"`
	PaymentMethod string          `json:"payment_method"` // "wechat" | "alipay"
	SkuIDs        []uint64        `json:"sku_ids"`         // subset of cart SKUs to checkout
	Remark        string          `json:"remark"`
}

func generateOrderNo() string {
	return fmt.Sprintf("%d%06d", time.Now().UnixMilli(), time.Now().Nanosecond()%1000000)
}

// CreateOrder validates cart items, deducts stock, and persists the order.
func CreateOrder(ctx context.Context, req CreateOrderReq) (*ordermodel.Order, error) {
	// 1. Load address
	var addr ordermodel.Address
	if err := db.DB.WithContext(ctx).Where("id = ? AND user_id = ?", req.AddressID, req.UserID).First(&addr).Error; err != nil {
		return nil, errors.New("收货地址不存在")
	}
	addrJSON, _ := json.Marshal(addr)

	// 2. Build order items from cart
	cartItems, err := GetCart(ctx, req.UserID)
	if err != nil || len(cartItems) == 0 {
		return nil, errors.New("购物车为空")
	}

	// Filter to requested SKU IDs
	skuSet := make(map[uint64]bool, len(req.SkuIDs))
	for _, id := range req.SkuIDs {
		skuSet[id] = true
	}

	var items []ordermodel.OrderItem
	var goodsAmount float64
	for _, ci := range cartItems {
		if len(req.SkuIDs) > 0 && !skuSet[ci.SkuID] {
			continue
		}
		if ci.Sku == nil || ci.Product == nil {
			continue
		}
		items = append(items, ordermodel.OrderItem{
			ProductID: ci.Product.ID,
			SkuID:     ci.SkuID,
			Title:     ci.Product.Title,
			Cover:     ci.Product.Cover,
			Attrs:     ci.Sku.Attrs,
			Price:     ci.Sku.Price,
			Qty:       ci.Qty,
		})
		goodsAmount += ci.Sku.Price * float64(ci.Qty)
	}
	if len(items) == 0 {
		return nil, errors.New("未选择有效商品")
	}

	// 3. Create order in transaction
	order := &ordermodel.Order{
		OrderNo:         generateOrderNo(),
		UserID:          req.UserID,
		Status:          ordermodel.OrderStatusPending,
		PaymentMethod:   req.PaymentMethod,
		GoodsAmount:     goodsAmount,
		TotalAmount:     goodsAmount,
		AddressSnapshot: addrJSON,
		Remark:          req.Remark,
	}

	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = order.ID
			// Deduct SKU stock
			res := tx.Model(&productmodel.ProductSku{}).
				Where("id = ? AND stock >= ?", items[i].SkuID, items[i].Qty).
				UpdateColumn("stock", gorm.Expr("stock - ?", items[i].Qty))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return fmt.Errorf("商品库存不足 (sku_id=%d)", items[i].SkuID)
			}
		}
		return tx.Create(&items).Error
	})
	if err != nil {
		return nil, err
	}

	// 4. Remove checked-out items from cart
	for _, item := range items {
		RemoveFromCart(ctx, req.UserID, item.SkuID)
	}

	return order, nil
}

// ListOrders returns paginated orders for a user.
func ListOrders(ctx context.Context, userID uint64, status int8, page, size int) ([]ordermodel.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Where("user_id = ?", userID)
	if status > 0 {
		tx = tx.Where("status = ?", status)
	}
	var total int64
	tx.Model(&ordermodel.Order{}).Count(&total)
	var list []ordermodel.Order
	err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

// ShipOrder marks an order as shipped (admin action).
func ShipOrder(ctx context.Context, orderID uint64, trackingNo string) error {
	return db.DB.WithContext(ctx).Model(&ordermodel.Order{}).
		Where("id = ? AND status = ?", orderID, ordermodel.OrderStatusPaid).
		Updates(map[string]any{"status": ordermodel.OrderStatusShipped}).Error
}

// CreateAddress saves a new address for a user.
func CreateAddress(ctx context.Context, addr *ordermodel.Address) error {
	return db.DB.WithContext(ctx).Create(addr).Error
}

// AdminListOrders returns orders with optional status filter for admin.
func AdminListOrders(ctx context.Context, status int8, page, size int) ([]ordermodel.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&ordermodel.Order{})
	if status > 0 {
		tx = tx.Where("status = ?", status)
	}
	var total int64
	tx.Count(&total)
	var list []ordermodel.Order
	err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}
