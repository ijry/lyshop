package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/marketing"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"gorm.io/gorm"
)

// CreateOrderReq is the request payload to create an order.
type CreateOrderReq struct {
	UserID        uint64   `json:"user_id"`
	AddressID     uint64   `json:"address_id"`
	PaymentMethod string   `json:"payment_method"`
	SkuIDs        []uint64 `json:"sku_ids"`
	CouponIDs     []uint64 `json:"coupon_ids"`
	PointsUse     int      `json:"points_use"`
	Remark        string   `json:"remark"`
}

type AmountBreakdown struct {
	GoodsAmount    float64 `json:"goods_amount"`
	DiscountAmount float64 `json:"discount_amount"`
	FreightAmount  float64 `json:"freight_amount"`
	PayableAmount  float64 `json:"payable_amount"`
}

type OrderView struct {
	ordermodel.Order
	Items            []OrderItemView            `json:"items"`
	AmountBreakdown  AmountBreakdown            `json:"amount_breakdown"`
	Shipments        []ordermodel.OrderShipment `json:"shipments,omitempty"`
	LatestShipment   *ordermodel.OrderShipment  `json:"latest_shipment,omitempty"`
	AfterSaleSummary *AfterSaleSummary          `json:"after_sale_summary,omitempty"`
}

type OrderItemView struct {
	ordermodel.OrderItem
	Review *ReviewView `json:"review,omitempty"`
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
	}
	if len(items) == 0 {
		return nil, errors.New("未选择有效商品")
	}

	// 3. Run pricing pipeline
	pricingItems := make([]marketing.OrderItem, len(items))
	for i, item := range items {
		pricingItems[i] = marketing.OrderItem{
			ProductID: item.ProductID, SkuID: item.SkuID,
			Title: item.Title, Price: item.Price, Qty: item.Qty,
		}
	}
	pCtx := &marketing.PriceContext{
		UserID:    req.UserID,
		Items:     pricingItems,
		CouponIDs: req.CouponIDs,
		PointsUse: req.PointsUse,
	}
	if err := marketing.Calculate(pCtx); err != nil {
		return nil, fmt.Errorf("价格计算失败: %w", err)
	}
	rulesJSON, _ := json.Marshal(pCtx.AppliedRules)

	// 4. Create order in transaction
	order := &ordermodel.Order{
		OrderNo:         generateOrderNo(),
		UserID:          req.UserID,
		Status:          ordermodel.OrderStatusPending,
		PaymentMethod:   req.PaymentMethod,
		GoodsAmount:     pCtx.GoodsAmount,
		DiscountAmount:  pCtx.ActivityDiscount + pCtx.FullReduceDiscount + pCtx.CouponDiscount + pCtx.PointsDiscount,
		TotalAmount:     pCtx.FinalAmount,
		AddressSnapshot: addrJSON,
		Remark:          req.Remark + string(rulesJSON),
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

func buildOrderViews(ctx context.Context, orders []ordermodel.Order) ([]OrderView, error) {
	if len(orders) == 0 {
		return []OrderView{}, nil
	}
	orderIDs := make([]uint64, 0, len(orders))
	for _, item := range orders {
		orderIDs = append(orderIDs, item.ID)
	}

	var orderItems []ordermodel.OrderItem
	if err := db.DB.WithContext(ctx).
		Where("order_id IN ?", orderIDs).
		Order("id asc").
		Find(&orderItems).Error; err != nil {
		return nil, err
	}

	itemMap := make(map[uint64][]ordermodel.OrderItem, len(orderIDs))
	itemIDs := make([]uint64, 0, len(orderItems))
	for _, item := range orderItems {
		itemMap[item.OrderID] = append(itemMap[item.OrderID], item)
		itemIDs = append(itemIDs, item.ID)
	}

	reviewMap := map[uint64]ReviewView{}
	if len(itemIDs) > 0 {
		var reviews []ordermodel.OrderReview
		if err := db.DB.WithContext(ctx).
			Where("order_item_id IN ?", itemIDs).
			Order("id asc").
			Find(&reviews).Error; err != nil {
			return nil, err
		}
		if len(reviews) > 0 {
			list, err := buildReviewView(ctx, reviews)
			if err != nil {
				return nil, err
			}
			for _, review := range list {
				reviewMap[review.OrderItemID] = review
			}
		}
	}

	shipmentsMap, err := listOrderShipmentsMap(ctx, orderIDs)
	if err != nil {
		return nil, err
	}
	afterSaleSummaryMap, err := buildAfterSaleSummaryMap(ctx, orderIDs)
	if err != nil {
		return nil, err
	}

	result := make([]OrderView, 0, len(orders))
	for _, item := range orders {
		items := make([]OrderItemView, 0, len(itemMap[item.ID]))
		for _, orderItem := range itemMap[item.ID] {
			var reviewPtr *ReviewView
			if review, ok := reviewMap[orderItem.ID]; ok {
				reviewCopy := review
				reviewPtr = &reviewCopy
			}
			items = append(items, OrderItemView{
				OrderItem: orderItem,
				Review:    reviewPtr,
			})
		}
		var latestShipment *ordermodel.OrderShipment
		orderShipments := shipmentsMap[item.ID]
		if len(orderShipments) > 0 {
			first := orderShipments[0]
			latestShipment = &first
		}
		result = append(result, OrderView{
			Order: item,
			Items: items,
			AmountBreakdown: AmountBreakdown{
				GoodsAmount:    item.GoodsAmount,
				DiscountAmount: item.DiscountAmount,
				FreightAmount:  item.FreightAmount,
				PayableAmount:  item.TotalAmount,
			},
			Shipments:        orderShipments,
			LatestShipment:   latestShipment,
			AfterSaleSummary: afterSaleSummaryMap[item.ID],
		})
	}
	return result, nil
}

// ListOrders returns paginated orders for a user.
func ListOrders(ctx context.Context, userID uint64, status int8, page, size int) ([]OrderView, int64, error) {
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
	var orders []ordermodel.Order
	if err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&orders).Error; err != nil {
		return nil, total, err
	}
	list, err := buildOrderViews(ctx, orders)
	return list, total, err
}

type ShipOrderReq struct {
	ShipType        string `json:"ship_type"`
	AfterSaleCaseID uint64 `json:"after_sale_case_id"`
	Company         string `json:"company"`
	TrackingNo      string `json:"tracking_no"`
	Remark          string `json:"remark"`
}

// ShipOrder marks an order as shipped or reshipped (admin action).
func ShipOrder(ctx context.Context, orderID uint64, req ShipOrderReq) error {
	req.TrackingNo = strings.TrimSpace(req.TrackingNo)
	if req.TrackingNo == "" {
		return errors.New("请填写快递单号")
	}
	shipType := strings.ToLower(strings.TrimSpace(req.ShipType))
	if shipType == "" {
		shipType = string(ordermodel.ShipmentBizTypeInitial)
	}
	if shipType != string(ordermodel.ShipmentBizTypeInitial) && shipType != string(ordermodel.ShipmentBizTypeReship) {
		return errors.New("不支持的发货类型")
	}
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order ordermodel.Order
		if err := tx.Where("id = ?", orderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("订单不存在")
			}
			return err
		}
		if shipType == string(ordermodel.ShipmentBizTypeInitial) && order.Status != ordermodel.OrderStatusPaid {
			return errors.New("当前状态不可发货")
		}
		if shipType == string(ordermodel.ShipmentBizTypeReship) {
			if req.AfterSaleCaseID == 0 {
				return errors.New("补发需关联售后单")
			}
			caseRow, err := lockAfterSaleCase(tx, req.AfterSaleCaseID)
			if err != nil {
				return err
			}
			if caseRow.OrderID != orderID {
				return errors.New("补发售后单与订单不匹配")
			}
			if caseRow.Status != string(ordermodel.AfterSaleStatusReshipPending) {
				return errors.New("当前售后状态不可补发")
			}
			if err := tx.Model(&ordermodel.AfterSaleCase{}).Where("id = ?", caseRow.ID).Update("status", string(ordermodel.AfterSaleStatusReshipped)).Error; err != nil {
				return err
			}
			if err := writeAfterSaleLogTx(tx, caseRow.ID, caseRow.Status, string(ordermodel.AfterSaleStatusReshipped), "reship", "admin", 0, "售后补发", map[string]any{
				"tracking_no": req.TrackingNo,
			}); err != nil {
				return err
			}
		}
		if _, err := createShipmentTx(tx, CreateShipmentReq{
			OrderID:         orderID,
			AfterSaleCaseID: req.AfterSaleCaseID,
			ShipType:        shipType,
			Direction:       string(ordermodel.ShipmentDirectionOutbound),
			Company:         req.Company,
			TrackingNo:      req.TrackingNo,
			Remark:          req.Remark,
			CreatedByType:   "admin",
			CreatedByID:     0,
		}); err != nil {
			return err
		}
		updates := map[string]any{"tracking_no": req.TrackingNo}
		if shipType == string(ordermodel.ShipmentBizTypeInitial) {
			updates["status"] = ordermodel.OrderStatusShipped
		}
		if err := tx.Model(&ordermodel.Order{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
			return err
		}
		if shipType == string(ordermodel.ShipmentBizTypeReship) {
			return refreshOrderStatusByAfterSaleTx(tx, orderID)
		}
		return nil
	})
}

// CreateAddress saves a new address for a user.
func CreateAddress(ctx context.Context, addr *ordermodel.Address) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&ordermodel.Address{}).Where("user_id = ?", addr.UserID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			addr.IsDefault = 1
		}
		if addr.IsDefault == 1 {
			if err := tx.Model(&ordermodel.Address{}).Where("user_id = ?", addr.UserID).Update("is_default", 0).Error; err != nil {
				return err
			}
		}
		return tx.Create(addr).Error
	})
}

func ListAddresses(ctx context.Context, userID uint64) ([]ordermodel.Address, error) {
	var list []ordermodel.Address
	err := db.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default desc, id desc").
		Find(&list).Error
	return list, err
}

func UpdateAddress(ctx context.Context, userID, id uint64, req ordermodel.Address) (*ordermodel.Address, error) {
	var addr ordermodel.Address
	if err := db.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&addr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("地址不存在")
		}
		return nil, err
	}

	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if req.IsDefault == 1 {
			if err := tx.Model(&ordermodel.Address{}).Where("user_id = ?", userID).Update("is_default", 0).Error; err != nil {
				return err
			}
		}

		updates := map[string]any{
			"name":       strings.TrimSpace(req.Name),
			"phone":      strings.TrimSpace(req.Phone),
			"province":   strings.TrimSpace(req.Province),
			"city":       strings.TrimSpace(req.City),
			"district":   strings.TrimSpace(req.District),
			"detail":     strings.TrimSpace(req.Detail),
			"is_default": req.IsDefault,
		}
		return tx.Model(&ordermodel.Address{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates).Error
	})
	if err != nil {
		return nil, err
	}
	if err := db.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&addr).Error; err != nil {
		return nil, err
	}
	return &addr, nil
}

func DeleteAddress(ctx context.Context, userID, id uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var target ordermodel.Address
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&target).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("地址不存在")
			}
			return err
		}
		if err := tx.Delete(&ordermodel.Address{}, target.ID).Error; err != nil {
			return err
		}
		if target.IsDefault == 1 {
			var another ordermodel.Address
			if err := tx.Where("user_id = ?", userID).Order("id desc").First(&another).Error; err == nil {
				if err := tx.Model(&ordermodel.Address{}).Where("id = ?", another.ID).Update("is_default", 1).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func PayOrder(ctx context.Context, userID, orderID uint64) error {
	now := time.Now()
	res := db.DB.WithContext(ctx).Model(&ordermodel.Order{}).
		Where("id = ? AND user_id = ? AND status = ?", orderID, userID, ordermodel.OrderStatusPending).
		Updates(map[string]any{
			"status":  ordermodel.OrderStatusPaid,
			"paid_at": &now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("订单不存在或当前状态不可支付")
	}
	return nil
}

func ReviewOrder(ctx context.Context, userID, orderID uint64, content string) error {
	var items []ordermodel.OrderItem
	if err := db.DB.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return err
	}
	reviewItems := make([]ReviewItemInput, 0, len(items))
	for _, item := range items {
		reviewItems = append(reviewItems, ReviewItemInput{
			OrderItemID:  item.ID,
			ProductScore: 5,
			Content:      strings.TrimSpace(content),
		})
	}
	return SubmitOrderReview(ctx, SubmitOrderReviewReq{
		OrderID:        orderID,
		UserID:         userID,
		Mode:           ReviewModeCreate,
		LogisticsScore: 5,
		Items:          reviewItems,
	})
}

func GetOrderDetail(ctx context.Context, userID, orderID uint64) (*OrderView, error) {
	var order ordermodel.Order
	if err := db.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}
	list, err := buildOrderViews(ctx, []ordermodel.Order{order})
	if err != nil {
		return nil, err
	}
	return &list[0], nil
}

func AdminGetOrderDetail(ctx context.Context, orderID uint64) (*OrderView, error) {
	var order ordermodel.Order
	if err := db.DB.WithContext(ctx).First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}
	list, err := buildOrderViews(ctx, []ordermodel.Order{order})
	if err != nil {
		return nil, err
	}
	return &list[0], nil
}

// AdminListOrders returns orders with optional status filter for admin.
func AdminListOrders(ctx context.Context, status int8, page, size int) ([]OrderView, int64, error) {
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
	var orders []ordermodel.Order
	if err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&orders).Error; err != nil {
		return nil, total, err
	}
	list, err := buildOrderViews(ctx, orders)
	return list, total, err
}
