package service

import (
	"context"
	"time"

	ordersvc "github.com/ijry/lyshop/plugins/order/service"
)

func StartPollingLoop() {
	if !loadBool("polling_enabled", false) {
		return
	}
	intervalSeconds := loadInt("polling_interval_seconds", 300)
	if intervalSeconds < 30 {
		intervalSeconds = 30
	}
	batchSize := loadInt("polling_batch_size", 100)
	if batchSize <= 0 {
		batchSize = 100
	}
	go func() {
		ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			_, _ = ordersvc.PollAndSyncShipments(context.Background(), batchSize)
		}
	}()
}
