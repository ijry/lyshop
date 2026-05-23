package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/ijry/lyshop/core/db"
	checkinmodel "github.com/ijry/lyshop/plugins/checkin/model"
	"gorm.io/gorm"
)

// Checkin performs the daily check-in for a user.
// Returns the awarded points and the consecutive day count.
func Checkin(ctx context.Context, userID uint64) (points int, consecutive int, err error) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Check if already checked in today
	var existing checkinmodel.CheckinLog
	if err := db.DB.WithContext(ctx).Where("user_id = ? AND date = ?", userID, today).First(&existing).Error; err == nil {
		return 0, existing.ConsecutiveDays, errors.New("今日已签到")
	}

	// Calculate consecutive days
	consecutive = 1
	var yesterdayLog checkinmodel.CheckinLog
	if err := db.DB.WithContext(ctx).Where("user_id = ? AND date = ?", userID, yesterday).First(&yesterdayLog).Error; err == nil {
		consecutive = yesterdayLog.ConsecutiveDays + 1
	}

	// Load rules and find matching points
	points = getPointsForDay(ctx, consecutive)

	// Create log
	log := &checkinmodel.CheckinLog{
		UserID:          userID,
		Date:            today,
		ConsecutiveDays: consecutive,
		Points:          points,
	}
	if err := db.DB.WithContext(ctx).Create(log).Error; err != nil {
		return 0, 0, err
	}

	// Add points to user (update user.points)
	db.DB.WithContext(ctx).Exec("UPDATE users SET points = points + ? WHERE id = ?", points, userID)

	return points, consecutive, nil
}

func getPointsForDay(ctx context.Context, day int) int {
	var rules []checkinmodel.CheckinRule
	db.DB.WithContext(ctx).Order("day asc").Find(&rules)

	if len(rules) == 0 {
		return 10 // default fallback
	}

	// Find the best matching rule
	sort.Slice(rules, func(i, j int) bool { return rules[i].Day < rules[j].Day })
	result := rules[0].Points // day=0 default
	for _, r := range rules {
		if r.Day == 0 {
			result = r.Points
		} else if r.Day == day {
			return r.Points
		} else if r.Day <= day {
			result = r.Points
		}
	}
	return result
}

// GetStatus returns the user's check-in status for today and this month.
func GetStatus(ctx context.Context, userID uint64) (map[string]any, error) {
	today := time.Now().Format("2006-01-02")
	monthStart := time.Now().Format("2006-01") + "-01"

	// Today checked?
	var todayLog checkinmodel.CheckinLog
	checkedToday := db.DB.WithContext(ctx).Where("user_id = ? AND date = ?", userID, today).First(&todayLog).Error == nil

	// Consecutive days
	consecutive := 0
	if checkedToday {
		consecutive = todayLog.ConsecutiveDays
	} else {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		var yLog checkinmodel.CheckinLog
		if db.DB.WithContext(ctx).Where("user_id = ? AND date = ?", userID, yesterday).First(&yLog).Error == nil {
			consecutive = yLog.ConsecutiveDays
		}
	}

	// This month's check-in dates
	var logs []checkinmodel.CheckinLog
	db.DB.WithContext(ctx).Where("user_id = ? AND date >= ?", userID, monthStart).Order("date asc").Find(&logs)
	dates := make([]string, len(logs))
	totalPoints := 0
	for i, l := range logs {
		dates[i] = l.Date
		totalPoints += l.Points
	}

	return map[string]any{
		"checked_today":    checkedToday,
		"consecutive_days": consecutive,
		"month_dates":      dates,
		"month_count":      len(dates),
		"month_points":     totalPoints,
	}, nil
}

// GetRules returns all check-in rules sorted by day.
func GetRules(ctx context.Context) ([]checkinmodel.CheckinRule, error) {
	var rules []checkinmodel.CheckinRule
	err := db.DB.WithContext(ctx).Order("day asc").Find(&rules).Error
	return rules, err
}

// SaveRules replaces all rules.
func SaveRules(ctx context.Context, rules []checkinmodel.CheckinRule) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Where("1=1").Delete(&checkinmodel.CheckinRule{})
		if len(rules) > 0 {
			return tx.Create(&rules).Error
		}
		return nil
	})
}

// AdminLogs returns check-in logs for admin view.
func AdminLogs(ctx context.Context, page, size int) ([]checkinmodel.CheckinLog, int64, error) {
	if page <= 0 { page = 1 }
	if size <= 0 || size > 100 { size = 20 }
	var total int64
	db.DB.Model(&checkinmodel.CheckinLog{}).Count(&total)
	var list []checkinmodel.CheckinLog
	err := db.DB.Order("id desc").Offset((page-1)*size).Limit(size).Find(&list).Error
	return list, total, err
}
