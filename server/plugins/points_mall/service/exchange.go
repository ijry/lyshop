package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/model"
	pmmodel "github.com/ijry/lyshop/plugins/points_mall/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ExchangeProduct 兑换积分商品
func ExchangeProduct(ctx context.Context, userID, productID uint64, qty int, addressSnapshot json.RawMessage) (*pmmodel.PointsExchange, error) {
	var exchange *pmmodel.PointsExchange

	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 锁定商品行
		var product pmmodel.PointsProduct
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&product, productID).Error; err != nil {
			return fmt.Errorf("商品不存在")
		}

		// 检查商品状态
		if product.Status != 1 {
			return fmt.Errorf("商品已下架")
		}

		// 2. 检查库存
		if product.Stock > 0 && product.Stock < qty {
			return fmt.Errorf("库存不足")
		}

		// 3. 检查用户积分
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return fmt.Errorf("用户不存在")
		}
		totalPoints := product.PointsPrice * qty
		if user.Points < totalPoints {
			return fmt.Errorf("积分不足")
		}

		// 4. 检查兑换限制
		if err := checkExchangeLimit(tx, userID, productID, qty, &product); err != nil {
			return err
		}

		// 5. 扣减库存
		if product.Stock > 0 {
			if err := tx.Model(&product).
				Update("stock", gorm.Expr("stock - ?", qty)).
				Update("sold_count", gorm.Expr("sold_count + ?", qty)).
				Error; err != nil {
				return err
			}
		} else {
			tx.Model(&product).Update("sold_count", gorm.Expr("sold_count + ?", qty))
		}

		// 6. 扣减积分
		if err := tx.Model(&user).
			Update("points", gorm.Expr("points - ?", totalPoints)).Error; err != nil {
			return err
		}

		// 7. 记录积分日志
		pointsLog := &pmmodel.PointsLog{
			UserID:    userID,
			Type:      3, // 兑换消耗
			Points:    -totalPoints,
			Remark:    fmt.Sprintf("兑换商品：%s", product.Title),
		}
		if err := tx.Create(pointsLog).Error; err != nil {
			return err
		}

		// 8. 创建兑换记录
		exchange = &pmmodel.PointsExchange{
			ExchangeNo:      generateExchangeNo(),
			UserID:          userID,
			ProductID:       productID,
			ProductTitle:    product.Title,
			ProductType:     product.Type,
			ProductCover:    product.Cover,
			PointsCost:      totalPoints,
			Qty:             qty,
			Status:          getInitialStatus(product.Type),
			AddressSnapshot: addressSnapshot,
			VirtualContent:  product.VirtualContent,
		}

		if err := tx.Create(exchange).Error; err != nil {
			return err
		}

		// 9. 特殊类型处理
		switch product.Type {
		case "coupon":
			return handleCouponExchange(tx, exchange, &product)
		case "virtual":
			// 虚拟商品立即完成
			now := time.Now()
			exchange.Status = "completed"
			exchange.CompletedAt = &now
			return tx.Save(exchange).Error
		}

		return nil
	})

	return exchange, err
}

// checkExchangeLimit 检查兑换限制
func checkExchangeLimit(tx *gorm.DB, userID, productID uint64, qty int, product *pmmodel.PointsProduct) error {
	// 检查每人限兑
	if product.LimitPerUser > 0 {
		var userExchangeCount int64
		tx.Model(&pmmodel.PointsExchange{}).
			Where("user_id = ? AND product_id = ? AND status != ?", userID, productID, "canceled").
			Select("COALESCE(SUM(qty), 0)").
			Scan(&userExchangeCount)

		if int(userExchangeCount)+qty > product.LimitPerUser {
			return fmt.Errorf("该商品每人限兑%d件", product.LimitPerUser)
		}
	}

	// 检查每日限兑
	if product.LimitPerDay > 0 {
		var todayExchangeCount int64
		tx.Model(&pmmodel.PointsExchange{}).
			Where("user_id = ? AND product_id = ? AND DATE(created_at) = CURDATE() AND status != ?",
				userID, productID, "canceled").
			Select("COALESCE(SUM(qty), 0)").
			Scan(&todayExchangeCount)

		if int(todayExchangeCount)+qty > product.LimitPerDay {
			return fmt.Errorf("该商品每日限兑%d件", product.LimitPerDay)
		}
	}

	return nil
}

// generateExchangeNo 生成兑换单号
func generateExchangeNo() string {
	return fmt.Sprintf("EX%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)
}

// getInitialStatus 获取初始状态
func getInitialStatus(productType string) string {
	switch productType {
	case "coupon", "virtual":
		return "completed" // 优惠券和虚拟商品立即完成
	default:
		return "pending_ship" // 实物商品待发货
	}
}

// handleCouponExchange 处理优惠券兑换
func handleCouponExchange(tx *gorm.DB, exchange *pmmodel.PointsExchange, product *pmmodel.PointsProduct) error {
	if product.CouponID == 0 {
		return fmt.Errorf("优惠券商品未配置关联优惠券")
	}

	// 调用 marketing 插件的优惠券发放服务
	// 创建优惠券用户记录
	couponUser := map[string]interface{}{
		"user_id":   exchange.UserID,
		"coupon_id": product.CouponID,
		"status":    1, // 未使用
		"source":    "points_exchange",
	}

	if err := tx.Table("coupon_users").Create(couponUser).Error; err != nil {
		return fmt.Errorf("发放优惠券失败: %w", err)
	}

	// 获取插入的ID
	var couponUserID uint64
	tx.Raw("SELECT LAST_INSERT_ID()").Scan(&couponUserID)
	exchange.CouponUserID = couponUserID

	now := time.Now()
	exchange.Status = "completed"
	exchange.CompletedAt = &now

	return tx.Save(exchange).Error
}

// ListExchanges 获取兑换记录列表
func ListExchanges(ctx context.Context, userID uint64, status string, page, size int) ([]pmmodel.PointsExchange, int64, error) {
	var exchanges []pmmodel.PointsExchange
	var total int64

	query := db.DB.WithContext(ctx).Model(&pmmodel.PointsExchange{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.Order("id DESC").Offset(offset).Limit(size).Find(&exchanges).Error

	return exchanges, total, err
}

// GetExchange 获取兑换详情
func GetExchange(ctx context.Context, id uint64) (*pmmodel.PointsExchange, error) {
	var exchange pmmodel.PointsExchange
	err := db.DB.WithContext(ctx).First(&exchange, id).Error
	return &exchange, err
}

// ShipExchange 发货
func ShipExchange(ctx context.Context, id uint64, trackingNo string) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exchange pmmodel.PointsExchange
		if err := tx.First(&exchange, id).Error; err != nil {
			return err
		}

		if exchange.Status != "pending_ship" {
			return fmt.Errorf("当前状态不允许发货")
		}

		now := time.Now()
		return tx.Model(&exchange).Updates(map[string]interface{}{
			"status":      "shipped",
			"tracking_no": trackingNo,
			"shipped_at":  &now,
		}).Error
	})
}

// CompleteExchange 完成兑换
func CompleteExchange(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exchange pmmodel.PointsExchange
		if err := tx.First(&exchange, id).Error; err != nil {
			return err
		}

		if exchange.Status != "shipped" && exchange.Status != "pending_ship" {
			return fmt.Errorf("当前状态不允许完成")
		}

		now := time.Now()
		return tx.Model(&exchange).Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": &now,
		}).Error
	})
}

// CancelExchange 取消兑换（退还积分）
func CancelExchange(ctx context.Context, id uint64, reason string) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exchange pmmodel.PointsExchange
		if err := tx.First(&exchange, id).Error; err != nil {
			return err
		}

		if exchange.Status == "completed" || exchange.Status == "canceled" {
			return fmt.Errorf("当前状态不允许取消")
		}

		// 退还积分
		if err := tx.Model(&model.User{}).
			Where("id = ?", exchange.UserID).
			Update("points", gorm.Expr("points + ?", exchange.PointsCost)).Error; err != nil {
			return err
		}

		// 记录积分日志
		pointsLog := &pmmodel.PointsLog{
			UserID:    exchange.UserID,
			Type:      5, // 管理员调整
			Points:    exchange.PointsCost,
			RelatedID: exchange.ID,
			Remark:    fmt.Sprintf("取消兑换退还积分：%s（原因：%s）", exchange.ProductTitle, reason),
		}
		if err := tx.Create(pointsLog).Error; err != nil {
			return err
		}

		// 恢复库存
		if err := tx.Model(&pmmodel.PointsProduct{}).
			Where("id = ? AND stock > 0", exchange.ProductID).
			Update("stock", gorm.Expr("stock + ?", exchange.Qty)).Error; err != nil {
			return err
		}

		// 更新兑换状态
		return tx.Model(&exchange).Updates(map[string]interface{}{
			"status": "canceled",
			"remark": reason,
		}).Error
	})
}
