package admin

import (
	"context"
	"time"

	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
)

type DashboardTrend struct {
	Date   string  `json:"date"`
	Orders int64   `json:"orders"`
	Sales  float64 `json:"sales"`
}

type DashboardCompare struct {
	RevenueYoY float64 `json:"revenue_yoy"`
	RevenueMoM float64 `json:"revenue_mom"`
	OrderYoY   float64 `json:"order_yoy"`
	OrderMoM   float64 `json:"order_mom"`
}

type DashboardStatusDistribution struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type DashboardHotProduct struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Cover   string `json:"cover"`
	SoldQty int64  `json:"sold_qty"`
}

type DashboardStockWarning struct {
	ProductID int64  `json:"product_id"`
	SkuID     int64  `json:"sku_id"`
	Title     string `json:"title"`
	Stock     int64  `json:"stock"`
	Threshold int64  `json:"threshold"`
}

type DashboardData struct {
	TodayOrders        int64                         `json:"today_orders"`
	TodaySales         float64                       `json:"today_sales"`
	TodayAvgPrice      float64                       `json:"today_avg_price"`
	PendingShip        int64                         `json:"pending_ship"`
	PendingAfterSale   int64                         `json:"pending_after_sale"`
	PendingRefunds     int64                         `json:"pending_refunds"`
	OnlineSessions     int64                         `json:"online_sessions"`
	UnreadMessage      int64                         `json:"unread_message"`
	StockWarning       int64                         `json:"stock_warning"`
	Compare            DashboardCompare              `json:"compare"`
	SalesTrend         []DashboardTrend              `json:"sales_trend"`
	StatusDistribution []DashboardStatusDistribution `json:"status_distribution"`
	HotProducts        []DashboardHotProduct         `json:"hot_products"`
	StockWarningList   []DashboardStockWarning       `json:"stock_warning_list"`
}

func GetDashboard(ctx context.Context) (*DashboardData, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)
	trendStart := todayStart.AddDate(0, 0, -29) // 30 days

	data := &DashboardData{
		SalesTrend:         make([]DashboardTrend, 0, 30),
		StatusDistribution: make([]DashboardStatusDistribution, 0),
		HotProducts:        make([]DashboardHotProduct, 0),
		StockWarningList:   make([]DashboardStockWarning, 0),
		Compare:            DashboardCompare{},
	}

	trendIndex := map[string]int{}
	for i := 29; i >= 0; i-- {
		day := todayStart.AddDate(0, 0, -i)
		date := day.Format("2006-01-02")
		trendIndex[date] = len(data.SalesTrend)
		data.SalesTrend = append(data.SalesTrend, DashboardTrend{
			Date:   date,
			Orders: 0,
			Sales:  0,
		})
	}

	if db.DB == nil {
		return data, nil
	}
	tx := db.DB.WithContext(ctx)

	if db.DB.Migrator().HasTable(&ordermodel.Order{}) {
		if err := tx.Model(&ordermodel.Order{}).
			Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
			Count(&data.TodayOrders).Error; err != nil {
			return nil, err
		}

		salesStatuses := []int8{
			ordermodel.OrderStatusPaid,
			ordermodel.OrderStatusShipped,
			ordermodel.OrderStatusCompleted,
			ordermodel.OrderStatusAfterSale,
		}
		if err := tx.Model(&ordermodel.Order{}).
			Where("created_at >= ? AND created_at < ? AND status IN ?", todayStart, todayEnd, salesStatuses).
			Select("COALESCE(SUM(total_amount), 0)").
			Scan(&data.TodaySales).Error; err != nil {
			return nil, err
		}

		if data.TodayOrders > 0 {
			data.TodayAvgPrice = data.TodaySales / float64(data.TodayOrders)
		}

		if err := tx.Model(&ordermodel.Order{}).
			Where("status = ?", ordermodel.OrderStatusPaid).
			Count(&data.PendingShip).Error; err != nil {
			return nil, err
		}

		var rows []struct {
			Day    string  `json:"day"`
			Orders int64   `json:"orders"`
			Sales  float64 `json:"sales"`
		}
		if err := tx.Model(&ordermodel.Order{}).
			Select("DATE(created_at) as day, COUNT(*) as orders, COALESCE(SUM(total_amount), 0) as sales").
			Where("created_at >= ? AND created_at < ? AND status IN ?", trendStart, todayEnd, salesStatuses).
			Group("DATE(created_at)").
			Scan(&rows).Error; err != nil {
			return nil, err
		}
		for _, row := range rows {
			if idx, ok := trendIndex[row.Day]; ok {
				data.SalesTrend[idx].Orders = row.Orders
				data.SalesTrend[idx].Sales = row.Sales
			}
		}
	}

	if db.DB.Migrator().HasTable(&ordermodel.AfterSaleCase{}) {
		closedStatuses := []string{
			string(ordermodel.AfterSaleStatusCompleted),
			string(ordermodel.AfterSaleStatusRejected),
			string(ordermodel.AfterSaleStatusClosed),
		}
		if err := tx.Model(&ordermodel.AfterSaleCase{}).
			Where("status NOT IN ?", closedStatuses).
			Count(&data.PendingAfterSale).Error; err != nil {
			return nil, err
		}
		data.PendingRefunds = data.PendingAfterSale
	}

	if db.DB.Migrator().HasTable(&immodel.ImSession{}) {
		if err := tx.Model(&immodel.ImSession{}).
			Where("status <> ?", immodel.SessionStatusClosed).
			Count(&data.OnlineSessions).Error; err != nil {
			return nil, err
		}
	}

	return data, nil
}
