package service

import (
	"context"

	ordersvc "github.com/ijry/lyshop/plugins/order/service"
)

func ListProductReviews(ctx context.Context, productID uint64, page, size int) (*ordersvc.ProductReviewList, error) {
	return ordersvc.ListProductReviews(ctx, productID, page, size)
}
