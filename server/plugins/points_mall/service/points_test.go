package service

import (
	"context"
	"testing"

	pmmodel "github.com/ijry/lyshop/plugins/points_mall/model"
)

func TestAddPoints(t *testing.T) {
	// 这是一个示例测试，实际使用时需要配置测试数据库
	ctx := context.Background()
	userID := uint64(1)
	points := 100
	logType := int8(5) // 管理员调整
	remark := "测试赠送积分"

	// 注意：这个测试需要数据库连接，实际运行前需要初始化数据库
	// 这里只是展示测试结构
	t.Skip("需要数据库连接")

	err := AddPoints(ctx, userID, points, logType, remark)
	if err != nil {
		t.Errorf("AddPoints failed: %v", err)
	}
}

func TestGetUserPoints(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	t.Skip("需要数据库连接")

	points, err := GetUserPoints(ctx, userID)
	if err != nil {
		t.Errorf("GetUserPoints failed: %v", err)
	}

	if points < 0 {
		t.Errorf("Expected non-negative points, got %d", points)
	}
}

func TestPointsLogModel(t *testing.T) {
	log := pmmodel.PointsLog{
		UserID:    1,
		Type:      1,
		Points:    100,
		RelatedID: 0,
		Remark:    "测试日志",
	}

	if log.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", log.UserID)
	}

	if log.Points != 100 {
		t.Errorf("Expected Points 100, got %d", log.Points)
	}
}

func TestPointsProductModel(t *testing.T) {
	product := pmmodel.PointsProduct{
		Title:       "测试商品",
		Type:        "physical",
		PointsPrice: 1000,
		Stock:       100,
		Status:      1,
	}

	if product.Title != "测试商品" {
		t.Errorf("Expected Title '测试商品', got '%s'", product.Title)
	}

	if product.PointsPrice != 1000 {
		t.Errorf("Expected PointsPrice 1000, got %d", product.PointsPrice)
	}
}

func TestPointsExchangeModel(t *testing.T) {
	exchange := pmmodel.PointsExchange{
		ExchangeNo:   "EX20260531000001",
		UserID:       1,
		ProductID:    1,
		ProductTitle: "测试商品",
		ProductType:  "physical",
		PointsCost:   1000,
		Qty:          1,
		Status:       "pending_ship",
	}

	if exchange.Status != "pending_ship" {
		t.Errorf("Expected Status 'pending_ship', got '%s'", exchange.Status)
	}

	if exchange.PointsCost != 1000 {
		t.Errorf("Expected PointsCost 1000, got %d", exchange.PointsCost)
	}
}
