package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ijry/lyshop/core/cache"
	"github.com/ijry/lyshop/core/db"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
)

func cartKey(userID uint64) string {
	return fmt.Sprintf("cart:%d", userID)
}

type CartItem struct {
	SkuID             uint64                   `json:"sku_id"`
	ActivityProductID uint64                   `json:"activity_product_id"`
	Qty               int                      `json:"qty"`
	Product           *productmodel.Product    `json:"product,omitempty"`
	Sku               *productmodel.ProductSku `json:"sku,omitempty"`
}

func formatCartField(skuID, activityProductID uint64) string {
	return fmt.Sprintf("%d:%d", skuID, activityProductID)
}

func parseCartField(field string) (uint64, uint64, error) {
	trimmed := strings.TrimSpace(field)
	if trimmed == "" {
		return 0, 0, fmt.Errorf("empty cart field")
	}
	if !strings.Contains(trimmed, ":") {
		// Backward compatibility: old key format was plain sku_id.
		skuID, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return 0, 0, err
		}
		return skuID, 0, nil
	}
	parts := strings.SplitN(trimmed, ":", 2)
	skuID, err := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	activityProductID, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return skuID, activityProductID, nil
}

func optionalActivityProductID(vals []uint64) uint64 {
	if len(vals) == 0 {
		return 0
	}
	return vals[0]
}

// AddToCart adds qty units of skuID to the user's Redis cart.
func AddToCart(ctx context.Context, userID, skuID uint64, qty int, activityProductID ...uint64) error {
	key := cartKey(userID)
	field := formatCartField(skuID, optionalActivityProductID(activityProductID))
	// Get existing qty
	existing, _ := cache.Client.HGet(ctx, key, field).Int()
	return cache.Client.HSet(ctx, key, field, existing+qty).Err()
}

// RemoveFromCart removes a SKU from the cart.
func RemoveFromCart(ctx context.Context, userID, skuID uint64, activityProductID ...uint64) error {
	field := formatCartField(skuID, optionalActivityProductID(activityProductID))
	return cache.Client.HDel(ctx, cartKey(userID), field).Err()
}

// GetCart returns all cart items with product/SKU details populated.
func GetCart(ctx context.Context, userID uint64) ([]CartItem, error) {
	fields, err := cache.Client.HGetAll(ctx, cartKey(userID)).Result()
	if err != nil {
		return nil, err
	}
	var items []CartItem
	for field, val := range fields {
		skuID, activityProductID, err := parseCartField(field)
		if err != nil {
			continue
		}
		qty, _ := strconv.Atoi(val)
		if qty <= 0 || skuID == 0 {
			continue
		}
		item := CartItem{SkuID: skuID, ActivityProductID: activityProductID, Qty: qty}
		var sku productmodel.ProductSku
		if err := db.DB.WithContext(ctx).First(&sku, skuID).Error; err == nil {
			item.Sku = &sku
			var prod productmodel.Product
			if err := db.DB.WithContext(ctx).First(&prod, sku.ProductID).Error; err == nil {
				item.Product = &prod
			}
		}
		items = append(items, item)
	}
	return items, nil
}

// ClearCart deletes the entire cart for a user.
func ClearCart(ctx context.Context, userID uint64) error {
	return cache.Client.Del(ctx, cartKey(userID)).Err()
}

// UpdateCartQty sets the quantity of a specific SKU (0 = remove).
func UpdateCartQty(ctx context.Context, userID, skuID uint64, qty int, activityProductID ...uint64) error {
	if qty <= 0 {
		return RemoveFromCart(ctx, userID, skuID, activityProductID...)
	}
	field := formatCartField(skuID, optionalActivityProductID(activityProductID))
	return cache.Client.HSet(ctx, cartKey(userID), field, qty).Err()
}

// skuIDsFromCart parses the sku IDs string (comma-separated) from cart fields.
func skuIDsFromCart(fields map[string]string) []string {
	ids := make([]string, 0, len(fields))
	for k := range fields {
		ids = append(ids, k)
	}
	return ids
}

// unused — kept for reference
var _ = strings.Join(skuIDsFromCart(nil), ",")
