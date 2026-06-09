package service

import (
	"context"
	"sync"

	"github.com/ijry/lyshop/core/inventory"
)

var (
	workerOnce sync.Once
	workerCtx  context.Context
)

func StartWorker() {
	if !inventory.IsAsyncExternalProvider() {
		return
	}
	workerOnce.Do(func() {
		workerCtx = context.Background()
		inventory.StartExternalWMSWorker(workerCtx, inventory.ExternalWorkerName())
	})
}
