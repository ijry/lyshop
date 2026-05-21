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
	SkuID   uint64                  `json:"sku_id"`
	Qty     int                     `json:"qty"`
	Product *productmodel.Product   `json:"product,omitempty"`
	Sku     *productmodel.ProductSku `json:"sku,omitempty"`
}

// AddToCart adds qty units of skuID to the user's Redis cart.
func AddToCart(ctx context.Context, userID, skuID uint64, qty int) error {
	key := cartKey(userID)
	field := strconv.FormatUint(skuID, 10)
	// Get existing qty
	existing, _ := cache.Client.HGet(ctx, key, field).Int()
	return cache.Client.HSet(ctx, key, field, existing+qty).Err()
}

// RemoveFromCart removes a SKU from the cart.
func RemoveFromCart(ctx context.Context, userID, skuID uint64) error {
	return cache.Client.HDel(ctx, cartKey(userID), strconv.FormatUint(skuID, 10)).Err()
}

// GetCart returns all cart items with product/SKU details populated.
func GetCart(ctx context.Context, userID uint64) ([]CartItem, error) {
	fields, err := cache.Client.HGetAll(ctx, cartKey(userID)).Result()
	if err != nil {
		return nil, err
	}
	var items []CartItem
	for field, val := range fields {
		skuID, _ := strconv.ParseUint(field, 10, 64)
		qty, _ := strconv.Atoi(val)
		if qty <= 0 || skuID == 0 {
			continue
		}
		item := CartItem{SkuID: skuID, Qty: qty}
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
func UpdateCartQty(ctx context.Context, userID, skuID uint64, qty int) error {
	if qty <= 0 {
		return RemoveFromCart(ctx, userID, skuID)
	}
	return cache.Client.HSet(ctx, cartKey(userID), strconv.FormatUint(skuID, 10), qty).Err()
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
