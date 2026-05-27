package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ActivityListQuery struct {
	Page int
	Size int
}

type ActivityProductListQuery struct {
	ActivityID uint64
	ProductID  uint64
	CategoryID uint64
	Keyword    string
	MinPrice   float64
	MaxPrice   float64
	SortBy     string
	SortOrder  string
	Page       int
	Size       int
}

type ManagedActivityProduct struct {
	mktmodel.ActivityProduct
	ProductTitle string                 `json:"product_title"`
	ProductCover string                 `json:"product_cover"`
	ProductPrice float64                `json:"product_price"`
	SkuAttrs     []productmodel.SkuAttr `json:"sku_attrs"`
	SkuPrice     float64                `json:"sku_price"`
	SkuStock     int                    `json:"sku_stock"`
}

type ActivityProductUpsert struct {
	ProductID       uint64  `json:"product_id"`
	SkuID           uint64  `json:"sku_id"`
	ActivityPrice   float64 `json:"activity_price"`
	StartPrice      float64 `json:"start_price"`
	FloorPrice      float64 `json:"floor_price"`
	LimitPerOrder   int     `json:"limit_per_order"`
	TotalStockLimit int     `json:"total_stock_limit"`
}

type FrontActivityProduct struct {
	ActivityProductID uint64     `json:"activity_product_id"`
	ActivityID        uint64     `json:"activity_id"`
	ActivityType      string     `json:"activity_type"`
	ActivityName      string     `json:"activity_name"`
	ActivityStartAt   *time.Time `json:"activity_start_at"`
	ActivityEndAt     *time.Time `json:"activity_end_at"`
	ProductID         uint64     `json:"product_id"`
	SkuID             uint64     `json:"sku_id"`
	Title             string     `json:"title"`
	Subtitle          string     `json:"subtitle"`
	Cover             string     `json:"cover"`
	CategoryID        uint64     `json:"category_id"`
	Sales             int        `json:"sales"`
	Stock             int        `json:"stock"`
	OriginPrice       float64    `json:"origin_price"`
	Price             float64    `json:"price"`
	ActivityPrice     float64    `json:"activity_price"`
	StartPrice        float64    `json:"start_price"`
	FloorPrice        float64    `json:"floor_price"`
	LimitPerOrder     int        `json:"limit_per_order"`
	TotalStockLimit   int        `json:"total_stock_limit"`
	SoldQty           int        `json:"sold_qty"`
}

type FrontActivityProductDetail struct {
	ActivityProductID uint64     `json:"activity_product_id"`
	ActivityID        uint64     `json:"activity_id"`
	ActivityType      string     `json:"activity_type"`
	ActivityName      string     `json:"activity_name"`
	ActivityStatus    int8       `json:"activity_status"`
	ActivityStartAt   *time.Time `json:"activity_start_at"`
	ActivityEndAt     *time.Time `json:"activity_end_at"`
	ProductID         uint64     `json:"product_id"`
	SkuID             uint64     `json:"sku_id"`
	Title             string     `json:"title"`
	Subtitle          string     `json:"subtitle"`
	Cover             string     `json:"cover"`
	CategoryID        uint64     `json:"category_id"`
	Sales             int        `json:"sales"`
	Stock             int        `json:"stock"`
	OriginPrice       float64    `json:"origin_price"`
	Price             float64    `json:"price"`
	ActivityPrice     float64    `json:"activity_price"`
	StartPrice        float64    `json:"start_price"`
	FloorPrice        float64    `json:"floor_price"`
	LimitPerOrder     int        `json:"limit_per_order"`
	TotalStockLimit   int        `json:"total_stock_limit"`
	SoldQty           int        `json:"sold_qty"`
}

var getFrontActivityProductDetailFn = GetFrontActivityProductDetail

func normalizePage(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	return page, size
}

func validActivityType(actType string) bool {
	switch actType {
	case mktmodel.ActivityTypeSeckill, mktmodel.ActivityTypeGroupBuy, mktmodel.ActivityTypeBargain:
		return true
	default:
		return false
	}
}

func normalizeActivityType(actType string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(actType))
	if normalized == "group-buy" {
		normalized = mktmodel.ActivityTypeGroupBuy
	}
	if !validActivityType(normalized) {
		return "", errors.New("invalid activity type")
	}
	return normalized, nil
}

func normalizeSort(sortBy, sortOrder string) (string, string) {
	sb := strings.TrimSpace(strings.ToLower(sortBy))
	so := strings.TrimSpace(strings.ToLower(sortOrder))
	if sb != "price" && sb != "sales" {
		sb = "price"
	}
	if so != "asc" && so != "desc" {
		so = "asc"
	}
	return sb, so
}

func checkActivityTimeConflict(ctx context.Context, actType string, startAt, endAt *time.Time, excludeID uint64) error {
	if startAt == nil || endAt == nil {
		return nil
	}
	if !startAt.Before(*endAt) {
		return errors.New("活动开始时间必须早于结束时间")
	}
	tx := db.DB.WithContext(ctx).Model(&mktmodel.Activity{}).
		Where("type = ?", actType).
		Where("start_at < ? AND end_at > ?", *endAt, *startAt)
	if excludeID > 0 {
		tx = tx.Where("id <> ?", excludeID)
	}
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("同类型活动时间段冲突")
	}
	return nil
}

func ListActivities(ctx context.Context, actType string, q ActivityListQuery) ([]mktmodel.Activity, int64, error) {
	normalizedType, err := normalizeActivityType(actType)
	if err != nil {
		return nil, 0, err
	}
	q.Page, q.Size = normalizePage(q.Page, q.Size)
	tx := db.DB.WithContext(ctx).Model(&mktmodel.Activity{}).Where("type = ?", normalizedType)
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []mktmodel.Activity
	if err := tx.Order("start_at desc, id desc").
		Offset((q.Page - 1) * q.Size).
		Limit(q.Size).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func CreateActivity(ctx context.Context, actType string, a *mktmodel.Activity) error {
	normalizedType, err := normalizeActivityType(actType)
	if err != nil {
		return err
	}
	a.Type = normalizedType
	if err := checkActivityTimeConflict(ctx, normalizedType, a.StartAt, a.EndAt, 0); err != nil {
		return err
	}
	return db.DB.WithContext(ctx).Create(a).Error
}

func UpdateActivity(ctx context.Context, activityID uint64, actType string, updates map[string]any) error {
	normalizedType, err := normalizeActivityType(actType)
	if err != nil {
		return err
	}
	var current mktmodel.Activity
	if err := db.DB.WithContext(ctx).Where("id = ? AND type = ?", activityID, normalizedType).First(&current).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("活动不存在")
		}
		return err
	}
	nextStart := current.StartAt
	nextEnd := current.EndAt
	if v, ok := updates["start_at"]; ok {
		nextStart = extractTimePtr(v, nextStart)
	}
	if v, ok := updates["end_at"]; ok {
		nextEnd = extractTimePtr(v, nextEnd)
	}
	if err := checkActivityTimeConflict(ctx, normalizedType, nextStart, nextEnd, activityID); err != nil {
		return err
	}
	return db.DB.WithContext(ctx).
		Model(&mktmodel.Activity{}).
		Where("id = ? AND type = ?", activityID, normalizedType).
		Updates(updates).Error
}

func ListManagedActivityProducts(ctx context.Context, actType string, q ActivityProductListQuery) ([]ManagedActivityProduct, int64, error) {
	normalizedType, err := normalizeActivityType(actType)
	if err != nil {
		return nil, 0, err
	}
	q.Page, q.Size = normalizePage(q.Page, q.Size)
	type row struct {
		mktmodel.ActivityProduct
		ProductTitle string
		ProductCover string
		ProductPrice float64
		SkuAttrsRaw  string
		SkuPrice     float64
		SkuStock     int
	}
	tx := db.DB.WithContext(ctx).
		Table("activity_products ap").
		Select(`
ap.id, ap.created_at, ap.updated_at, ap.activity_id, ap.product_id, ap.sku_id, ap.activity_price,
ap.start_price, ap.floor_price, ap.limit_per_order, ap.total_stock_limit, ap.sold_qty, ap.activity_stock,
p.title AS product_title, p.cover AS product_cover, p.price AS product_price,
ps.attrs AS sku_attrs_raw, ps.price AS sku_price, ps.stock AS sku_stock`).
		Joins("JOIN activities a ON a.id = ap.activity_id").
		Joins("JOIN products p ON p.id = ap.product_id").
		Joins("LEFT JOIN product_skus ps ON ps.id = ap.sku_id").
		Where("a.type = ?", normalizedType)
	if q.ActivityID > 0 {
		tx = tx.Where("ap.activity_id = ?", q.ActivityID)
	}
	if q.ProductID > 0 {
		tx = tx.Where("ap.product_id = ?", q.ProductID)
	}
	if q.CategoryID > 0 {
		tx = tx.Where("p.category_id = ?", q.CategoryID)
	}
	if strings.TrimSpace(q.Keyword) != "" {
		tx = tx.Where("p.title LIKE ?", "%"+strings.TrimSpace(q.Keyword)+"%")
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []row
	if err := tx.Order("ap.id desc").
		Offset((q.Page - 1) * q.Size).
		Limit(q.Size).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	result := make([]ManagedActivityProduct, 0, len(rows))
	for _, item := range rows {
		result = append(result, ManagedActivityProduct{
			ActivityProduct: item.ActivityProduct,
			ProductTitle:    item.ProductTitle,
			ProductCover:    item.ProductCover,
			ProductPrice:    item.ProductPrice,
			SkuAttrs:        decodeSkuAttrs(item.SkuAttrsRaw),
			SkuPrice:        item.SkuPrice,
			SkuStock:        item.SkuStock,
		})
	}
	return result, total, nil
}

func UpsertActivityProducts(ctx context.Context, activityID uint64, actType string, rows []ActivityProductUpsert) error {
	normalizedType, err := normalizeActivityType(actType)
	if err != nil {
		return err
	}
	var activity mktmodel.Activity
	if err := db.DB.WithContext(ctx).Where("id = ? AND type = ?", activityID, normalizedType).First(&activity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("活动不存在")
		}
		return err
	}
	if len(rows) == 0 {
		return db.DB.WithContext(ctx).Where("activity_id = ?", activityID).Delete(&mktmodel.ActivityProduct{}).Error
	}
	seen := make(map[uint64]struct{}, len(rows))
	for _, item := range rows {
		if item.ProductID == 0 || item.SkuID == 0 {
			return errors.New("商品SKU不能为空")
		}
		if _, ok := seen[item.SkuID]; ok {
			return errors.New("同一活动内SKU不可重复")
		}
		seen[item.SkuID] = struct{}{}
		if item.LimitPerOrder < 0 || item.TotalStockLimit < 0 {
			return errors.New("限购与活动库存必须大于等于0")
		}
		switch normalizedType {
		case mktmodel.ActivityTypeSeckill, mktmodel.ActivityTypeGroupBuy:
			if item.ActivityPrice <= 0 {
				return errors.New("活动价必须大于0")
			}
		case mktmodel.ActivityTypeBargain:
			if item.StartPrice <= 0 || item.FloorPrice <= 0 {
				return errors.New("砍价起始价和最低价必须大于0")
			}
			if item.FloorPrice > item.StartPrice {
				return errors.New("砍价最低价不能高于起始价")
			}
		}
	}
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var oldRows []mktmodel.ActivityProduct
		if err := tx.Where("activity_id = ?", activityID).Find(&oldRows).Error; err != nil {
			return err
		}
		oldSold := make(map[uint64]int, len(oldRows))
		for _, item := range oldRows {
			oldSold[item.SkuID] = item.SoldQty
		}
		if err := tx.Where("activity_id = ?", activityID).Delete(&mktmodel.ActivityProduct{}).Error; err != nil {
			return err
		}
		newRows := make([]mktmodel.ActivityProduct, 0, len(rows))
		for _, item := range rows {
			if item.TotalStockLimit > 0 && oldSold[item.SkuID] > item.TotalStockLimit {
				return errors.New("活动库存上限不能小于已售数量")
			}
			row := mktmodel.ActivityProduct{
				ActivityID:      activityID,
				ProductID:       item.ProductID,
				SkuID:           item.SkuID,
				ActivityPrice:   roundPrice(item.ActivityPrice),
				StartPrice:      roundPrice(item.StartPrice),
				FloorPrice:      roundPrice(item.FloorPrice),
				LimitPerOrder:   item.LimitPerOrder,
				TotalStockLimit: item.TotalStockLimit,
				SoldQty:         oldSold[item.SkuID],
			}
			newRows = append(newRows, row)
		}
		return tx.Create(&newRows).Error
	})
}

func ListFrontActivityProducts(ctx context.Context, actType string, q ActivityProductListQuery) ([]FrontActivityProduct, int64, error) {
	normalizedType, err := normalizeActivityType(actType)
	if err != nil {
		return nil, 0, err
	}
	q.Page, q.Size = normalizePage(q.Page, q.Size)
	q.SortBy, q.SortOrder = normalizeSort(q.SortBy, q.SortOrder)
	type joined struct {
		ActivityProductID uint64
		ActivityID        uint64
		ActivityType      string
		ActivityName      string
		ActivityStartAt   *time.Time
		ActivityEndAt     *time.Time
		ProductID         uint64
		SkuID             uint64
		Title             string
		Subtitle          string
		Cover             string
		CategoryID        uint64
		Sales             int
		Stock             int
		ProductPrice      float64
		ActivityPrice     float64
		StartPrice        float64
		FloorPrice        float64
		LimitPerOrder     int
		TotalStockLimit   int
		SoldQty           int
		SkuPrice          float64
	}
	now := time.Now()
	tx := db.DB.WithContext(ctx).
		Table("activity_products ap").
		Select(`
ap.id AS activity_product_id, ap.activity_id, a.type AS activity_type, a.name AS activity_name, a.start_at AS activity_start_at, a.end_at AS activity_end_at,
ap.product_id, ap.sku_id,
p.title, p.subtitle, p.cover, p.category_id, p.sales, p.stock,
p.price AS product_price, ap.activity_price, ap.start_price, ap.floor_price, ap.limit_per_order, ap.total_stock_limit, ap.sold_qty,
ps.price AS sku_price`).
		Joins("JOIN activities a ON a.id = ap.activity_id").
		Joins("JOIN products p ON p.id = ap.product_id").
		Joins("LEFT JOIN product_skus ps ON ps.id = ap.sku_id").
		Where("a.type = ? AND a.status = 1 AND a.start_at <= ? AND a.end_at >= ? AND p.status = 1", normalizedType, now, now)
	if q.CategoryID > 0 {
		tx = tx.Where("p.category_id = ?", q.CategoryID)
	}
	if strings.TrimSpace(q.Keyword) != "" {
		tx = tx.Where("p.title LIKE ?", "%"+strings.TrimSpace(q.Keyword)+"%")
	}
	var rows []joined
	if err := tx.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	result := make([]FrontActivityProduct, 0, len(rows))
	for _, item := range rows {
		originPrice := item.ProductPrice
		if item.SkuPrice > 0 {
			originPrice = item.SkuPrice
		}
		price := originPrice
		switch normalizedType {
		case mktmodel.ActivityTypeSeckill, mktmodel.ActivityTypeGroupBuy:
			if item.ActivityPrice > 0 {
				price = item.ActivityPrice
			}
		case mktmodel.ActivityTypeBargain:
			if item.StartPrice > 0 {
				price = item.StartPrice
			}
		}
		if q.MinPrice > 0 && price < q.MinPrice {
			continue
		}
		if q.MaxPrice > 0 && price > q.MaxPrice {
			continue
		}
		if item.TotalStockLimit > 0 && item.SoldQty >= item.TotalStockLimit {
			continue
		}
		result = append(result, FrontActivityProduct{
			ActivityProductID: item.ActivityProductID,
			ActivityID:        item.ActivityID,
			ActivityType:      item.ActivityType,
			ActivityName:      item.ActivityName,
			ActivityStartAt:   item.ActivityStartAt,
			ActivityEndAt:     item.ActivityEndAt,
			ProductID:         item.ProductID,
			SkuID:             item.SkuID,
			Title:             item.Title,
			Subtitle:          item.Subtitle,
			Cover:             item.Cover,
			CategoryID:        item.CategoryID,
			Sales:             item.Sales,
			Stock:             item.Stock,
			OriginPrice:       roundPrice(originPrice),
			Price:             roundPrice(price),
			ActivityPrice:     roundPrice(item.ActivityPrice),
			StartPrice:        roundPrice(item.StartPrice),
			FloorPrice:        roundPrice(item.FloorPrice),
			LimitPerOrder:     item.LimitPerOrder,
			TotalStockLimit:   item.TotalStockLimit,
			SoldQty:           item.SoldQty,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if q.SortBy == "sales" {
			if q.SortOrder == "desc" {
				if result[i].Sales == result[j].Sales {
					return result[i].ProductID > result[j].ProductID
				}
				return result[i].Sales > result[j].Sales
			}
			if result[i].Sales == result[j].Sales {
				return result[i].ProductID < result[j].ProductID
			}
			return result[i].Sales < result[j].Sales
		}
		if q.SortOrder == "desc" {
			if result[i].Price == result[j].Price {
				return result[i].ProductID > result[j].ProductID
			}
			return result[i].Price > result[j].Price
		}
		if result[i].Price == result[j].Price {
			return result[i].ProductID < result[j].ProductID
		}
		return result[i].Price < result[j].Price
	})
	total := int64(len(result))
	offset := (q.Page - 1) * q.Size
	if offset >= len(result) {
		return []FrontActivityProduct{}, total, nil
	}
	end := offset + q.Size
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], total, nil
}

func GetFrontActivityProductDetail(ctx context.Context, activityProductID uint64) (*FrontActivityProductDetail, error) {
	if activityProductID == 0 {
		return nil, errors.New("活动商品不存在")
	}
	type joined struct {
		ActivityProductID uint64
		ActivityID        uint64
		ActivityType      string
		ActivityName      string
		ActivityStatus    int8
		ActivityStartAt   *time.Time
		ActivityEndAt     *time.Time
		ProductID         uint64
		SkuID             uint64
		Title             string
		Subtitle          string
		Cover             string
		CategoryID        uint64
		Sales             int
		Stock             int
		ProductPrice      float64
		ActivityPrice     float64
		StartPrice        float64
		FloorPrice        float64
		LimitPerOrder     int
		TotalStockLimit   int
		SoldQty           int
		SkuPrice          float64
	}
	var row joined
	err := db.DB.WithContext(ctx).
		Table("activity_products ap").
		Select(`
ap.id AS activity_product_id, ap.activity_id, a.type AS activity_type, a.name AS activity_name, a.status AS activity_status, a.start_at AS activity_start_at, a.end_at AS activity_end_at,
ap.product_id, ap.sku_id,
p.title, p.subtitle, p.cover, p.category_id, p.sales, p.stock,
p.price AS product_price, ap.activity_price, ap.start_price, ap.floor_price, ap.limit_per_order, ap.total_stock_limit, ap.sold_qty,
ps.price AS sku_price`).
		Joins("JOIN activities a ON a.id = ap.activity_id").
		Joins("JOIN products p ON p.id = ap.product_id").
		Joins("LEFT JOIN product_skus ps ON ps.id = ap.sku_id").
		Where("ap.id = ? AND p.status = 1", activityProductID).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("活动商品不存在")
		}
		return nil, err
	}
	originPrice := row.ProductPrice
	if row.SkuPrice > 0 {
		originPrice = row.SkuPrice
	}
	price := originPrice
	switch row.ActivityType {
	case mktmodel.ActivityTypeSeckill, mktmodel.ActivityTypeGroupBuy:
		if row.ActivityPrice > 0 {
			price = row.ActivityPrice
		}
	case mktmodel.ActivityTypeBargain:
		if row.StartPrice > 0 {
			price = row.StartPrice
		}
	}
	return &FrontActivityProductDetail{
		ActivityProductID: row.ActivityProductID,
		ActivityID:        row.ActivityID,
		ActivityType:      row.ActivityType,
		ActivityName:      row.ActivityName,
		ActivityStatus:    row.ActivityStatus,
		ActivityStartAt:   row.ActivityStartAt,
		ActivityEndAt:     row.ActivityEndAt,
		ProductID:         row.ProductID,
		SkuID:             row.SkuID,
		Title:             row.Title,
		Subtitle:          row.Subtitle,
		Cover:             row.Cover,
		CategoryID:        row.CategoryID,
		Sales:             row.Sales,
		Stock:             row.Stock,
		OriginPrice:       roundPrice(originPrice),
		Price:             roundPrice(price),
		ActivityPrice:     roundPrice(row.ActivityPrice),
		StartPrice:        roundPrice(row.StartPrice),
		FloorPrice:        roundPrice(row.FloorPrice),
		LimitPerOrder:     row.LimitPerOrder,
		TotalStockLimit:   row.TotalStockLimit,
		SoldQty:           row.SoldQty,
	}, nil
}

func ValidateActivityProductSource(ctx context.Context, activityProductID, skuID, productID uint64) (*FrontActivityProductDetail, error) {
	detail, err := getFrontActivityProductDetailFn(ctx, activityProductID)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, errors.New("活动商品不存在")
	}
	if skuID > 0 && detail.SkuID > 0 && detail.SkuID != skuID {
		return nil, errors.New("活动商品SKU不匹配")
	}
	if productID > 0 && detail.ProductID != productID {
		return nil, errors.New("活动商品与商品不匹配")
	}
	if detail.ActivityStatus != 1 {
		return nil, errors.New("活动已失效")
	}
	now := time.Now()
	if detail.ActivityStartAt != nil && now.Before(*detail.ActivityStartAt) {
		return nil, errors.New("活动未开始")
	}
	if detail.ActivityEndAt != nil && now.After(*detail.ActivityEndAt) {
		return nil, errors.New("活动已结束")
	}
	if detail.TotalStockLimit > 0 && detail.SoldQty >= detail.TotalStockLimit {
		return nil, errors.New("活动库存不足")
	}
	return detail, nil
}

func IncreaseSoldQtyTx(tx *gorm.DB, skuID uint64, qty int) error {
	if skuID == 0 || qty <= 0 {
		return nil
	}
	now := time.Now()
	var rows []mktmodel.ActivityProduct
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Joins("JOIN activities a ON a.id = activity_products.activity_id").
		Where("activity_products.sku_id = ? AND a.status = 1 AND a.start_at <= ? AND a.end_at >= ?", skuID, now, now).
		Find(&rows).Error; err != nil {
		return err
	}
	for _, item := range rows {
		if item.LimitPerOrder > 0 && qty > item.LimitPerOrder {
			return errors.New("超过活动单笔限购")
		}
		if item.TotalStockLimit > 0 && item.SoldQty+qty > item.TotalStockLimit {
			return errors.New("活动库存不足")
		}
		if err := tx.Model(&mktmodel.ActivityProduct{}).
			Where("id = ?", item.ID).
			UpdateColumn("sold_qty", gorm.Expr("sold_qty + ?", qty)).Error; err != nil {
			return err
		}
	}
	return nil
}

func roundPrice(value float64) float64 {
	if value <= 0 {
		return 0
	}
	return math.Round(value*100) / 100
}

func decodeSkuAttrs(raw string) []productmodel.SkuAttr {
	if strings.TrimSpace(raw) == "" {
		return []productmodel.SkuAttr{}
	}
	var attrs []productmodel.SkuAttr
	if err := json.Unmarshal([]byte(raw), &attrs); err != nil {
		return []productmodel.SkuAttr{}
	}
	return attrs
}

func extractTimePtr(input any, fallback *time.Time) *time.Time {
	switch val := input.(type) {
	case *time.Time:
		return val
	case time.Time:
		t := val
		return &t
	default:
		return fallback
	}
}
