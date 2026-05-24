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

type DashboardData struct {
	TodayOrders    int64            `json:"today_orders"`
	TodaySales     float64          `json:"today_sales"`
	PendingRefunds int64            `json:"pending_refunds"`
	OnlineSessions int64            `json:"online_sessions"`
	SalesTrend     []DashboardTrend `json:"sales_trend"`
}

func GetDashboard(ctx context.Context) (*DashboardData, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)
	trendStart := todayStart.AddDate(0, 0, -6)

	data := &DashboardData{
		SalesTrend: make([]DashboardTrend, 0, 7),
	}

	trendIndex := map[string]int{}
	for i := 6; i >= 0; i-- {
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
			Count(&data.PendingRefunds).Error; err != nil {
			return nil, err
		}
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
