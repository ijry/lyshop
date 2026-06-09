package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	inventorycore "github.com/ijry/lyshop/core/inventory"
	"github.com/ijry/lyshop/core/marketing"
	marketingsvc "github.com/ijry/lyshop/plugins/marketing/service"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	vipsvc "github.com/ijry/lyshop/plugins/vip/service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateOrderReq is the request payload to create an order.
type CreateOrderReq struct {
	UserID        uint64               `json:"user_id"`
	AddressID     uint64               `json:"address_id"`
	PaymentMethod string               `json:"payment_method"`
	Items         []CreateOrderItemReq `json:"items"`
	SkuIDs        []uint64             `json:"sku_ids"`
	CouponIDs     []uint64             `json:"coupon_ids"`
	PointsUse     int                  `json:"points_use"`
	Remark        string               `json:"remark"`
}

type CreateOrderItemReq struct {
	SkuID             uint64 `json:"sku_id"`
	ActivityProductID uint64 `json:"activity_product_id"`
}

type AmountBreakdown struct {
	GoodsAmount    float64 `json:"goods_amount"`
	DiscountAmount float64 `json:"discount_amount"`
	FreightAmount  float64 `json:"freight_amount"`
	PayableAmount  float64 `json:"payable_amount"`
}

type OrderView struct {
	ordermodel.Order
	Items            []OrderItemView     `json:"items"`
	AmountBreakdown  AmountBreakdown     `json:"amount_breakdown"`
	Shipments        []OrderShipmentView `json:"shipments,omitempty"`
	LatestShipment   *OrderShipmentView  `json:"latest_shipment,omitempty"`
	AfterSaleSummary *AfterSaleSummary   `json:"after_sale_summary,omitempty"`
}

type OrderItemView struct {
	ordermodel.OrderItem
	Review *ReviewView `json:"review,omitempty"`
}

type OrderShipmentView struct {
	ordermodel.OrderShipment
	DeliveryTypeLabel    string `json:"delivery_type_label,omitempty"`
	DirectionLabel       string `json:"direction_label,omitempty"`
	BizTypeLabel         string `json:"biz_type_label,omitempty"`
	LogisticsStatusLabel string `json:"logistics_status_label,omitempty"`
}

var validateActivityProductSourceFn = marketingsvc.ValidateActivityProductSource
var increaseSoldQtyByActivityProductTxFn = marketingsvc.IncreaseSoldQtyByActivityProductTx
var getInventoryProviderFn = inventorycore.CurrentProvider

const (
	orderReservationBizType = "order"
	orderReserveExpire      = 15 * time.Minute
)

func cartSelectionKey(skuID, activityProductID uint64) string {
	return fmt.Sprintf("%d:%d", skuID, activityProductID)
}

func buildOrderItemsFromCart(ctx context.Context, cartItems []CartItem, req CreateOrderReq) ([]ordermodel.OrderItem, error) {
	selected := make(map[string]struct{})
	selectedBySKU := make(map[uint64]struct{})
	useItems := len(req.Items) > 0
	if useItems {
		for _, item := range req.Items {
			if item.SkuID == 0 {
				continue
			}
			selected[cartSelectionKey(item.SkuID, item.ActivityProductID)] = struct{}{}
		}
	} else {
		for _, id := range req.SkuIDs {
			if id == 0 {
				continue
			}
			selectedBySKU[id] = struct{}{}
		}
	}

	items := make([]ordermodel.OrderItem, 0, len(cartItems))
	for _, ci := range cartItems {
		if ci.Sku == nil || ci.Product == nil {
			continue
		}
		if useItems {
			if _, ok := selected[cartSelectionKey(ci.SkuID, ci.ActivityProductID)]; !ok {
				continue
			}
		} else if len(selectedBySKU) > 0 {
			if _, ok := selectedBySKU[ci.SkuID]; !ok {
				continue
			}
		}

		item := ordermodel.OrderItem{
			ProductID: ci.Product.ID,
			SkuID:     ci.SkuID,
			Title:     ci.Product.Title,
			Cover:     ci.Product.Cover,
			Attrs:     ci.Sku.Attrs,
			Price:     ci.Sku.Price,
			Qty:       ci.Qty,
		}
		if ci.ActivityProductID > 0 {
			detail, err := validateActivityProductSourceFn(ctx, ci.ActivityProductID, ci.SkuID, ci.Product.ID)
			if err != nil {
				return nil, err
			}
			item.ActivityProductID = detail.ActivityProductID
			item.ActivityID = detail.ActivityID
			item.ActivityType = detail.ActivityType
			item.ActivityTitle = detail.ActivityName
			if detail.Price > 0 {
				item.Price = detail.Price
			}
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, errors.New("order.err.noValidProducts")
	}
	return items, nil
}

func buildOrderShipmentViews(c *gin.Context, rows []ordermodel.OrderShipment) []OrderShipmentView {
	if len(rows) == 0 {
		return []OrderShipmentView{}
	}
	result := make([]OrderShipmentView, 0, len(rows))
	for _, row := range rows {
		result = append(result, OrderShipmentView{
			OrderShipment:        row,
			DeliveryTypeLabel:    deliveryTypeLabel(c, row.DeliveryType),
			DirectionLabel:       shipmentDirectionLabel(c, row.Direction),
			BizTypeLabel:         shipmentBizTypeLabel(c, row.BizType),
			LogisticsStatusLabel: shipmentStatusLabel(c, row.LogisticsStatus),
		})
	}
	return result
}

func generateOrderNo() string {
	return fmt.Sprintf("%d%06d", time.Now().UnixMilli(), time.Now().Nanosecond()%1000000)
}

func reserveOrderInventory(tx *gorm.DB, orderNo string, items []ordermodel.OrderItem) error {
	provider, err := getInventoryProviderFn()
	if err != nil {
		return err
	}
	reserveItems := make([]inventorycore.ReserveItem, 0, len(items))
	for _, item := range items {
		reserveItems = append(reserveItems, inventorycore.ReserveItem{
			SkuID: item.SkuID,
			Qty:   item.Qty,
		})
	}
	expireAt := time.Now().Add(orderReserveExpire)
	return provider.ReserveTx(tx, inventorycore.ReserveInput{
		BizType:   orderReservationBizType,
		BizNo:     orderNo,
		Items:     reserveItems,
		ExpiredAt: &expireAt,
	})
}

func confirmOrderInventory(tx *gorm.DB, orderNo string) error {
	provider, err := getInventoryProviderFn()
	if err != nil {
		return err
	}
	return provider.ConfirmTx(tx, orderReservationBizType, orderNo)
}

func releaseOrderInventory(tx *gorm.DB, orderNo, reason string) error {
	provider, err := getInventoryProviderFn()
	if err != nil {
		return err
	}
	return provider.ReleaseTx(tx, orderReservationBizType, orderNo, reason)
}

// CreateOrder validates cart items, deducts stock, and persists the order.
func CreateOrder(ctx context.Context, req CreateOrderReq) (*ordermodel.Order, error) {
	// 1. Load address
	var addr ordermodel.Address
	if err := db.DB.WithContext(ctx).Where("id = ? AND user_id = ?", req.AddressID, req.UserID).First(&addr).Error; err != nil {
		return nil, errors.New("order.err.addressNotFound")
	}
	addrJSON, _ := json.Marshal(addr)

	// 2. Build order items from cart
	cartItems, err := GetCart(ctx, req.UserID)
	if err != nil || len(cartItems) == 0 {
		return nil, errors.New("order.err.cartEmpty")
	}

	items, err := buildOrderItemsFromCart(ctx, cartItems, req)
	if err != nil {
		return nil, err
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
	if asset, err := vipsvc.GetActiveAsset(ctx, req.UserID, time.Now()); err == nil && asset != nil {
		pCtx.IsVIP = true
		pCtx.VIPLevelID = asset.CurrentLevelID
	}
	if err := marketing.Calculate(pCtx); err != nil {
		return nil, fmt.Errorf("order.err.priceCalcFailed: %w", err)
	}
	rulesJSON, _ := json.Marshal(pCtx.AppliedRules)

	// 4. Create order in transaction
	order := &ordermodel.Order{
		OrderNo:         generateOrderNo(),
		UserID:          req.UserID,
		Status:          ordermodel.OrderStatusPending,
		InventoryStatus: inventorycore.InventoryStatusNone,
		PaymentMethod:   req.PaymentMethod,
		GoodsAmount:     pCtx.GoodsAmount,
		DiscountAmount:  pCtx.ActivityDiscount + pCtx.VipDiscount + pCtx.FullReduceDiscount + pCtx.CouponDiscount + pCtx.PointsDiscount,
		TotalAmount:     pCtx.FinalAmount,
		AddressSnapshot: addrJSON,
		Remark:          req.Remark + string(rulesJSON),
	}

	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if err := reserveOrderInventory(tx, order.OrderNo, items); err != nil {
			return err
		}
		if err := tx.Model(&ordermodel.Order{}).Where("id = ?", order.ID).Update("inventory_status", inventorycore.InventoryStatusReserved).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = order.ID
			if items[i].ActivityProductID > 0 {
				if err := increaseSoldQtyByActivityProductTxFn(tx, items[i].ActivityProductID, items[i].Qty); err != nil {
					return err
				}
			}
		}
		return tx.Create(&items).Error
	})
	if err != nil {
		return nil, err
	}

	// 4. Remove checked-out items from cart
	for _, item := range items {
		RemoveFromCart(ctx, req.UserID, item.SkuID, item.ActivityProductID)
	}

	return order, nil
}

func buildOrderViews(c *gin.Context, ctx context.Context, orders []ordermodel.Order) ([]OrderView, error) {
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
	afterSaleSummaryMap, err := buildAfterSaleSummaryMap(c, ctx, orderIDs)
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
		var latestShipment *OrderShipmentView
		orderShipments := buildOrderShipmentViews(c, shipmentsMap[item.ID])
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
func ListOrders(c *gin.Context, userID uint64, status int8, page, size int) ([]OrderView, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 20
	}
	ctx := c.Request.Context()
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
	list, err := buildOrderViews(c, ctx, orders)
	return list, total, err
}

type ShipOrderReq struct {
	DeliveryType    string `json:"delivery_type"`
	ShipType        string `json:"ship_type"`
	AfterSaleCaseID uint64 `json:"after_sale_case_id"`
	Company         string `json:"company"`
	TrackingNo      string `json:"tracking_no"`
	RiderName       string `json:"rider_name"`
	RiderPhone      string `json:"rider_phone"`
	Remark          string `json:"remark"`
}

// ShipOrder marks an order as shipped or reshipped (admin action).
func ShipOrder(ctx context.Context, orderID uint64, req ShipOrderReq) error {
	deliveryType := normalizeDeliveryType(req.DeliveryType)
	if deliveryType == string(ordermodel.DeliveryTypeLocal) {
		if strings.TrimSpace(req.RiderName) == "" {
			return errors.New("delivery.err.riderNameRequired")
		}
		if strings.TrimSpace(req.RiderPhone) == "" {
			return errors.New("delivery.err.riderPhoneRequired")
		}
	} else {
		req.TrackingNo = strings.TrimSpace(req.TrackingNo)
		if req.TrackingNo == "" {
			return errors.New("delivery.err.trackingRequired")
		}
	}
	shipType := strings.ToLower(strings.TrimSpace(req.ShipType))
	if shipType == "" {
		shipType = string(ordermodel.ShipmentBizTypeInitial)
	}
	if shipType != string(ordermodel.ShipmentBizTypeInitial) && shipType != string(ordermodel.ShipmentBizTypeReship) {
		return errors.New("order.err.unsupportedDelivery")
	}
	var createdShipmentID uint64
	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order ordermodel.Order
		if err := tx.Where("id = ?", orderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("order.err.notFound")
			}
			return err
		}
		if shipType == string(ordermodel.ShipmentBizTypeInitial) && order.Status != ordermodel.OrderStatusPaid {
			return errors.New("order.err.cannotShip")
		}
		if shipType == string(ordermodel.ShipmentBizTypeReship) {
			if req.AfterSaleCaseID == 0 {
				return errors.New("order.err.reshipNeedCase")
			}
			caseRow, err := lockAfterSaleCase(tx, req.AfterSaleCaseID)
			if err != nil {
				return err
			}
			if caseRow.OrderID != orderID {
				return errors.New("order.err.reshipCaseMismatch")
			}
			if caseRow.Status != string(ordermodel.AfterSaleStatusReshipPending) {
				return errors.New("order.err.reshipStatusInvalid")
			}
			if err := tx.Model(&ordermodel.AfterSaleCase{}).Where("id = ?", caseRow.ID).Update("status", string(ordermodel.AfterSaleStatusReshipped)).Error; err != nil {
				return err
			}
			if err := writeAfterSaleLogTx(tx, caseRow.ID, caseRow.Status, string(ordermodel.AfterSaleStatusReshipped), "reship", "admin", 0, "afterSale.log.reship", map[string]any{
				"tracking_no": req.TrackingNo,
			}); err != nil {
				return err
			}
		}
		shipmentRow, err := createShipmentTx(tx, CreateShipmentReq{
			OrderID:         orderID,
			AfterSaleCaseID: req.AfterSaleCaseID,
			DeliveryType:    deliveryType,
			ShipType:        shipType,
			Direction:       string(ordermodel.ShipmentDirectionOutbound),
			Company:         req.Company,
			TrackingNo:      req.TrackingNo,
			RiderName:       req.RiderName,
			RiderPhone:      req.RiderPhone,
			Remark:          req.Remark,
			CreatedByType:   "admin",
			CreatedByID:     0,
		})
		if err != nil {
			return err
		}
		createdShipmentID = shipmentRow.ID
		updates := map[string]any{}
		if req.TrackingNo != "" {
			updates["tracking_no"] = req.TrackingNo
		}
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
	if err != nil {
		return err
	}
	if createdShipmentID > 0 && deliveryType != string(ordermodel.DeliveryTypeLocal) {
		go func(id uint64) {
			_ = SyncShipmentTracks(context.Background(), id, SyncShipmentReq{Manual: false})
		}(createdShipmentID)
	}
	return nil
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
			return nil, errors.New("address not found")
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
				return errors.New("address not found")
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
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order ordermodel.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", orderID, userID).
			First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("order.err.cannotPay")
			}
			return err
		}
		if order.Status != ordermodel.OrderStatusPending {
			return errors.New("order.err.cannotPay")
		}
		if err := confirmOrderInventory(tx, order.OrderNo); err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&ordermodel.Order{}).Where("id = ?", order.ID).Updates(map[string]any{
			"status":           ordermodel.OrderStatusPaid,
			"inventory_status": inventorycore.InventoryStatusConfirmed,
			"paid_at":          &now,
		}).Error; err != nil {
			return err
		}

		return vipsvc.GrantGrowthForPaidOrderTx(tx, userID, order.ID, order.TotalAmount)
	})
}

func CancelOrder(ctx context.Context, userID, orderID uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order ordermodel.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", orderID, userID).
			First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("order.err.cannotCancel")
			}
			return err
		}
		if order.Status != ordermodel.OrderStatusPending {
			return errors.New("order.err.cannotCancel")
		}
		if err := releaseOrderInventory(tx, order.OrderNo, "user_cancel"); err != nil {
			return err
		}
		return tx.Model(&ordermodel.Order{}).Where("id = ?", order.ID).Updates(map[string]any{
			"status":           ordermodel.OrderStatusCanceled,
			"inventory_status": inventorycore.InventoryStatusReleased,
		}).Error
	})
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

func GetOrderDetail(c *gin.Context, userID, orderID uint64) (*OrderView, error) {
	ctx := c.Request.Context()
	var order ordermodel.Order
	if err := db.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order.err.notFound")
		}
		return nil, err
	}
	list, err := buildOrderViews(c, ctx, []ordermodel.Order{order})
	if err != nil {
		return nil, err
	}
	return &list[0], nil
}

func AdminGetOrderDetail(c *gin.Context, orderID uint64) (*OrderView, error) {
	ctx := c.Request.Context()
	var order ordermodel.Order
	if err := db.DB.WithContext(ctx).First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order.err.notFound")
		}
		return nil, err
	}
	list, err := buildOrderViews(c, ctx, []ordermodel.Order{order})
	if err != nil {
		return nil, err
	}
	return &list[0], nil
}

// AdminListOrders returns orders with optional status filter for admin.
func AdminListOrders(c *gin.Context, status int8, page, size int) ([]OrderView, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	ctx := c.Request.Context()
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
	list, err := buildOrderViews(c, ctx, orders)
	return list, total, err
}
