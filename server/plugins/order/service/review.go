package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ijry/lyshop/core/db"
	adminmodel "github.com/ijry/lyshop/model"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"gorm.io/gorm"
)

type ReviewMode string

const (
	ReviewModeCreate ReviewMode = "create"
	ReviewModeEdit   ReviewMode = "edit"
	ReviewModeAppend ReviewMode = "append"
)

type ReviewItemInput struct {
	OrderItemID  uint64   `json:"order_item_id"`
	ProductScore int8     `json:"product_score"`
	Content      string   `json:"content"`
	Images       []string `json:"images"`
}

type SubmitOrderReviewReq struct {
	OrderID        uint64            `json:"order_id"`
	UserID         uint64            `json:"user_id"`
	Mode           ReviewMode        `json:"mode"`
	LogisticsScore int8              `json:"logistics_score"`
	Items          []ReviewItemInput `json:"items"`
	AppendContent  string            `json:"append_content"`
	AppendImages   []string          `json:"append_images"`
}

type ReviewAppendView struct {
	ordermodel.OrderReviewAppend
	Images []string `json:"images"`
}

type ReviewReplyView struct {
	ordermodel.OrderReviewReply
}

type ReviewView struct {
	ordermodel.OrderReview
	OrderNo      string                `json:"order_no"`
	OrderItem    *ordermodel.OrderItem `json:"order_item,omitempty"`
	Product      *productmodel.Product `json:"product,omitempty"`
	Images       []string              `json:"images"`
	Appends      []ReviewAppendView    `json:"appends"`
	AdminReply   *ReviewReplyView      `json:"admin_reply,omitempty"`
	UserNickname string                `json:"user_nickname,omitempty"`
	UserAvatar   string                `json:"user_avatar,omitempty"`
}

type ReviewOption struct {
	OrderItemID    uint64   `json:"order_item_id"`
	ReviewID       uint64   `json:"review_id"`
	HasReview      bool     `json:"has_review"`
	ProductID      uint64   `json:"product_id"`
	ProductTitle   string   `json:"product_title"`
	ProductCover   string   `json:"product_cover"`
	ProductScore   int8     `json:"product_score"`
	LogisticsScore int8     `json:"logistics_score"`
	Content        string   `json:"content"`
	Images         []string `json:"images"`
}

type OrderReviewMeta struct {
	OrderID        uint64         `json:"order_id"`
	OrderNo        string         `json:"order_no"`
	LogisticsScore int8           `json:"logistics_score"`
	CanCreate      bool           `json:"can_create"`
	CanEdit        bool           `json:"can_edit"`
	CanAppend      bool           `json:"can_append"`
	Options        []ReviewOption `json:"options"`
}

type ReviewSummary struct {
	AvgProductScore   float64 `json:"avg_product_score"`
	AvgLogisticsScore float64 `json:"avg_logistics_score"`
	Total             int64   `json:"total"`
}

type ProductReviewList struct {
	Summary ReviewSummary `json:"summary"`
	List    []ReviewView  `json:"list"`
	Total   int64         `json:"total"`
	Page    int           `json:"page"`
	Size    int           `json:"size"`
}

type AdminReviewList struct {
	List  []ReviewView `json:"list"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

type AdminReviewReplyReq struct {
	AdminID uint64 `json:"admin_id"`
	Content string `json:"content"`
}

func clampScore(score int8) int8 {
	if score <= 0 {
		return 5
	}
	if score > 5 {
		return 5
	}
	return score
}

func normalizeImageURLs(urls []string) []string {
	result := make([]string, 0, len(urls))
	seen := make(map[string]struct{}, len(urls))
	for _, raw := range urls {
		u := strings.TrimSpace(raw)
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		result = append(result, u)
		if len(result) >= 9 {
			break
		}
	}
	return result
}

func encodeImageURLs(urls []string) string {
	normalized := normalizeImageURLs(urls)
	if len(normalized) == 0 {
		return "[]"
	}
	buf, err := json.Marshal(normalized)
	if err != nil {
		return "[]"
	}
	return string(buf)
}

func decodeImageURLs(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	var urls []string
	if err := json.Unmarshal([]byte(raw), &urls); err != nil {
		return []string{}
	}
	return normalizeImageURLs(urls)
}

func buildReviewView(ctx context.Context, reviews []ordermodel.OrderReview) ([]ReviewView, error) {
	if len(reviews) == 0 {
		return []ReviewView{}, nil
	}

	reviewIDs := make([]uint64, 0, len(reviews))
	orderIDs := make([]uint64, 0, len(reviews))
	orderItemIDs := make([]uint64, 0, len(reviews))
	productIDs := make([]uint64, 0, len(reviews))
	userIDs := make([]uint64, 0, len(reviews))
	for _, r := range reviews {
		reviewIDs = append(reviewIDs, r.ID)
		orderIDs = append(orderIDs, r.OrderID)
		orderItemIDs = append(orderItemIDs, r.OrderItemID)
		productIDs = append(productIDs, r.ProductID)
		userIDs = append(userIDs, r.UserID)
	}

	var items []ordermodel.OrderItem
	if err := db.DB.WithContext(ctx).Where("id IN ?", orderItemIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	itemMap := make(map[uint64]ordermodel.OrderItem, len(items))
	for _, item := range items {
		itemMap[item.ID] = item
	}

	var products []productmodel.Product
	if err := db.DB.WithContext(ctx).Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		return nil, err
	}
	productMap := make(map[uint64]productmodel.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	var users []adminmodel.User
	if err := db.DB.WithContext(ctx).Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	userMap := make(map[uint64]adminmodel.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	var appends []ordermodel.OrderReviewAppend
	if err := db.DB.WithContext(ctx).Where("review_id IN ?", reviewIDs).Order("id asc").Find(&appends).Error; err != nil {
		return nil, err
	}
	appendMap := make(map[uint64][]ReviewAppendView, len(reviewIDs))
	for _, app := range appends {
		appendMap[app.ReviewID] = append(appendMap[app.ReviewID], ReviewAppendView{
			OrderReviewAppend: app,
			Images:            decodeImageURLs(app.ImagesJSON),
		})
	}

	var replies []ordermodel.OrderReviewReply
	if err := db.DB.WithContext(ctx).Where("review_id IN ?", reviewIDs).Find(&replies).Error; err != nil {
		return nil, err
	}
	replyMap := make(map[uint64]ReviewReplyView, len(reviewIDs))
	for _, reply := range replies {
		replyMap[reply.ReviewID] = ReviewReplyView{OrderReviewReply: reply}
	}

	var orders []ordermodel.Order
	if err := db.DB.WithContext(ctx).Where("id IN ?", orderIDs).Find(&orders).Error; err != nil {
		return nil, err
	}
	orderMap := make(map[uint64]ordermodel.Order, len(orders))
	for _, o := range orders {
		orderMap[o.ID] = o
	}

	result := make([]ReviewView, 0, len(reviews))
	for _, r := range reviews {
		view := ReviewView{
			OrderReview: r,
			OrderNo:     orderMap[r.OrderID].OrderNo,
			Images:      decodeImageURLs(r.ImagesJSON),
			Appends:     appendMap[r.ID],
		}
		if item, ok := itemMap[r.OrderItemID]; ok {
			itemCopy := item
			view.OrderItem = &itemCopy
		}
		if product, ok := productMap[r.ProductID]; ok {
			productCopy := product
			view.Product = &productCopy
		}
		if user, ok := userMap[r.UserID]; ok {
			view.UserNickname = user.Nickname
			view.UserAvatar = user.Avatar
		}
		if reply, ok := replyMap[r.ID]; ok {
			replyCopy := reply
			view.AdminReply = &replyCopy
		}
		result = append(result, view)
	}
	return result, nil
}

func loadRootReview(tx *gorm.DB, orderID, orderItemID, userID uint64) (*ordermodel.OrderReview, error) {
	var review ordermodel.OrderReview
	err := tx.
		Where("order_id = ? AND order_item_id = ? AND user_id = ?", orderID, orderItemID, userID).
		First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

func SubmitOrderReview(ctx context.Context, req SubmitOrderReviewReq) error {
	if req.UserID == 0 || req.OrderID == 0 {
		return errors.New("参数错误")
	}
	if req.Mode == "" {
		req.Mode = ReviewModeCreate
	}
	if len(req.Items) == 0 && req.Mode != ReviewModeAppend {
		return errors.New("评价内容不能为空")
	}
	if len(req.Items) == 0 && req.Mode == ReviewModeAppend {
		return errors.New("请先选择要追加的评价商品")
	}

	var order ordermodel.Order
	if err := db.DB.WithContext(ctx).First(&order, req.OrderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}
	if order.UserID != req.UserID {
		return errors.New("订单不存在")
	}
	if order.Status < ordermodel.OrderStatusShipped {
		return errors.New("订单未完成，暂不可评价")
	}

	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		switch req.Mode {
		case ReviewModeCreate:
			for _, item := range req.Items {
				if item.OrderItemID == 0 {
					continue
				}
				root, err := loadRootReview(tx, req.OrderID, item.OrderItemID, req.UserID)
				if err != nil {
					return err
				}
				if root != nil {
					return fmt.Errorf("商品 %d 已评价", item.OrderItemID)
				}
				var orderItem ordermodel.OrderItem
				if err := tx.Where("id = ? AND order_id = ?", item.OrderItemID, req.OrderID).First(&orderItem).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("订单商品不存在")
					}
					return err
				}
				review := &ordermodel.OrderReview{
					OrderID:        req.OrderID,
					OrderItemID:    item.OrderItemID,
					ProductID:      orderItem.ProductID,
					UserID:         req.UserID,
					MerchantID:     order.MerchantID,
					ProductScore:   clampScore(item.ProductScore),
					LogisticsScore: clampScore(req.LogisticsScore),
					Content:        strings.TrimSpace(item.Content),
					ImagesJSON:     encodeImageURLs(item.Images),
				}
				if err := tx.Create(review).Error; err != nil {
					return err
				}
			}
		case ReviewModeEdit:
			for _, item := range req.Items {
				if item.OrderItemID == 0 {
					continue
				}
				var review ordermodel.OrderReview
				if err := tx.Where("order_id = ? AND order_item_id = ? AND user_id = ?", req.OrderID, item.OrderItemID, req.UserID).First(&review).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("评价不存在")
					}
					return err
				}
				updates := map[string]any{
					"product_score":   clampScore(item.ProductScore),
					"logistics_score": clampScore(req.LogisticsScore),
					"content":         strings.TrimSpace(item.Content),
					"images_json":     encodeImageURLs(item.Images),
					"edited_times":    gorm.Expr("edited_times + ?", 1),
				}
				if err := tx.Model(&ordermodel.OrderReview{}).Where("id = ?", review.ID).Updates(updates).Error; err != nil {
					return err
				}
			}
		case ReviewModeAppend:
			appendContent := strings.TrimSpace(req.AppendContent)
			appendImages := normalizeImageURLs(req.AppendImages)
			if appendContent == "" && len(appendImages) == 0 {
				return errors.New("追评内容或图片不能为空")
			}
			createdAny := false
			for _, item := range req.Items {
				if item.OrderItemID == 0 {
					continue
				}
				root, err := loadRootReview(tx, req.OrderID, item.OrderItemID, req.UserID)
				if err != nil {
					return err
				}
				if root == nil {
					return errors.New("只能对根评价追加")
				}
				appendRow := &ordermodel.OrderReviewAppend{
					ReviewID: root.ID,
					UserID:   req.UserID,
					Content:  appendContent,
					ImagesJSON: encodeImageURLs(appendImages),
				}
				if err := tx.Create(appendRow).Error; err != nil {
					return err
				}
				createdAny = true
			}
			if !createdAny {
				return errors.New("未找到可追加的评价")
			}
		default:
			return errors.New("不支持的评价模式")
		}

		if req.Mode != ReviewModeAppend {
			var totalItems int64
			if err := tx.Model(&ordermodel.OrderItem{}).Where("order_id = ?", req.OrderID).Count(&totalItems).Error; err != nil {
				return err
			}
			var reviewedItems int64
			if err := tx.Model(&ordermodel.OrderReview{}).Where("order_id = ? AND user_id = ?", req.OrderID, req.UserID).Count(&reviewedItems).Error; err != nil {
				return err
			}
			if totalItems > 0 && reviewedItems >= totalItems {
				if err := tx.Model(&ordermodel.Order{}).Where("id = ? AND user_id = ?", req.OrderID, req.UserID).
					Update("status", ordermodel.OrderStatusCompleted).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func ListUserOrderReviews(ctx context.Context, userID, orderID uint64) ([]ReviewView, error) {
	var reviews []ordermodel.OrderReview
	if err := db.DB.WithContext(ctx).
		Where("order_id = ? AND user_id = ?", orderID, userID).
		Order("id asc").
		Find(&reviews).Error; err != nil {
		return nil, err
	}
	return buildReviewView(ctx, reviews)
}

func GetOrderReviewMeta(ctx context.Context, userID, orderID uint64) (*OrderReviewMeta, error) {
	var order ordermodel.Order
	if err := db.DB.WithContext(ctx).Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}
	var items []ordermodel.OrderItem
	if err := db.DB.WithContext(ctx).Where("order_id = ?", orderID).Order("id asc").Find(&items).Error; err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("订单商品不存在")
	}
	itemIDs := make([]uint64, 0, len(items))
	for _, item := range items {
		itemIDs = append(itemIDs, item.ID)
	}
	var reviews []ordermodel.OrderReview
	if err := db.DB.WithContext(ctx).
		Where("order_id = ? AND user_id = ? AND order_item_id IN ?", orderID, userID, itemIDs).
		Order("id asc").
		Find(&reviews).Error; err != nil {
		return nil, err
	}
	reviewByItemID := make(map[uint64]ordermodel.OrderReview, len(reviews))
	for _, review := range reviews {
		reviewByItemID[review.OrderItemID] = review
	}
	logisticsScore := int8(5)
	if len(reviews) > 0 {
		logisticsScore = reviews[0].LogisticsScore
	}
	options := make([]ReviewOption, 0, len(items))
	for _, item := range items {
		opt := ReviewOption{
			OrderItemID:    item.ID,
			ProductID:      item.ProductID,
			ProductTitle:   item.Title,
			ProductCover:   item.Cover,
			ProductScore:   5,
			LogisticsScore: logisticsScore,
		}
		if review, ok := reviewByItemID[item.ID]; ok {
			opt.ReviewID = review.ID
			opt.HasReview = true
			opt.ProductScore = review.ProductScore
			opt.LogisticsScore = review.LogisticsScore
			opt.Content = review.Content
			opt.Images = decodeImageURLs(review.ImagesJSON)
		}
		options = append(options, opt)
	}
	canCreate := len(reviews) < len(items)
	canEdit := len(reviews) > 0
	canAppend := len(reviews) > 0
	return &OrderReviewMeta{
		OrderID:        orderID,
		OrderNo:        order.OrderNo,
		LogisticsScore: logisticsScore,
		CanCreate:      canCreate,
		CanEdit:        canEdit,
		CanAppend:      canAppend,
		Options:        options,
	}, nil
}

func ListProductReviews(ctx context.Context, productID uint64, page, size int) (*ProductReviewList, error) {
	if productID == 0 {
		return nil, errors.New("商品不存在")
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	tx := db.DB.WithContext(ctx).Model(&ordermodel.OrderReview{}).Where("product_id = ?", productID)
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}
	var reviews []ordermodel.OrderReview
	if err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&reviews).Error; err != nil {
		return nil, err
	}
	list, err := buildReviewView(ctx, reviews)
	if err != nil {
		return nil, err
	}
	var avgProduct, avgLogistics float64
	// compute summary explicitly to avoid DB-specific scan issues
	if total > 0 {
		var agg struct {
			AvgProductScore   float64 `gorm:"column:avg_product_score"`
			AvgLogisticsScore float64 `gorm:"column:avg_logistics_score"`
		}
		if err := db.DB.WithContext(ctx).Model(&ordermodel.OrderReview{}).
			Select("AVG(product_score) as avg_product_score, AVG(logistics_score) as avg_logistics_score").
			Where("product_id = ?", productID).
			Scan(&agg).Error; err != nil {
			return nil, err
		}
		avgProduct = agg.AvgProductScore
		avgLogistics = agg.AvgLogisticsScore
	}
	return &ProductReviewList{
		Summary: ReviewSummary{
			AvgProductScore:   avgProduct,
			AvgLogisticsScore: avgLogistics,
			Total:             total,
		},
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

func AdminListReviews(ctx context.Context, productID uint64, keyword string, page, size int) (*AdminReviewList, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&ordermodel.OrderReview{})
	if productID > 0 {
		tx = tx.Where("product_id = ?", productID)
	}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		tx = tx.Where("content LIKE ?", "%"+keyword+"%")
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}
	var reviews []ordermodel.OrderReview
	if err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&reviews).Error; err != nil {
		return nil, err
	}
	list, err := buildReviewView(ctx, reviews)
	if err != nil {
		return nil, err
	}
	return &AdminReviewList{List: list, Total: total, Page: page, Size: size}, nil
}

func AdminGetReview(ctx context.Context, id uint64) (*ReviewView, error) {
	var review ordermodel.OrderReview
	if err := db.DB.WithContext(ctx).First(&review, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评价不存在")
		}
		return nil, err
	}
	list, err := buildReviewView(ctx, []ordermodel.OrderReview{review})
	if err != nil {
		return nil, err
	}
	return &list[0], nil
}

func AdminUpsertReply(ctx context.Context, reviewID uint64, req AdminReviewReplyReq) error {
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("回复内容不能为空")
	}
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var review ordermodel.OrderReview
		if err := tx.First(&review, reviewID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("评价不存在")
			}
			return err
		}
		var reply ordermodel.OrderReviewReply
		err := tx.Where("review_id = ?", reviewID).First(&reply).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(&ordermodel.OrderReviewReply{
				ReviewID: reviewID,
				AdminID:  req.AdminID,
				Content:  strings.TrimSpace(req.Content),
			}).Error
		}
		return tx.Model(&ordermodel.OrderReviewReply{}).Where("id = ?", reply.ID).
			Updates(map[string]any{"admin_id": req.AdminID, "content": strings.TrimSpace(req.Content)}).Error
	})
}
