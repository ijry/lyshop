package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/ijry/lyshop/config"
	"github.com/ijry/lyshop/core/db"
)

func StartExternalWMSWorker(ctx context.Context, workerName string) {
	go func() {
		interval := time.Duration(config.Global.ExternalWMS.WorkerIntervalSec) * time.Second
		if interval <= 0 {
			interval = 5 * time.Second
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		processor := &externalProvider{adapter: NewExternalAdapter()}
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runExternalWMSTick(processor, workerName)
			}
		}
	}()
}

func runExternalWMSTick(processor AsyncTaskProcessor, workerName string) {
	now := time.Now()
	for {
		task, err := ClaimDueTask(db.DB, workerName, now)
		if err != nil || task == nil {
			return
		}
		if err := ProcessTask(db.DB, processor, task, workerName, now); err != nil {
			return
		}
	}
}

func ExternalWorkerName() string {
	return fmt.Sprintf("external-wms-%d", time.Now().UnixNano())
}
