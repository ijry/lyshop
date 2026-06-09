package inventory

import (
	"context"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type Provider interface {
	Name() string
	ReserveTx(tx *gorm.DB, in ReserveInput) error
	ConfirmTx(tx *gorm.DB, bizType, bizNo string) error
	ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error
	DeductTx(tx *gorm.DB, in DeductInput) error
	RestoreTx(tx *gorm.DB, in RestoreInput) error
	SyncSkuTx(tx *gorm.DB, in SyncSkuInput) error
	GetSellableStock(ctx context.Context, skuIDs []uint64) ([]SellableStock, error)
}

var (
	regMu     sync.RWMutex
	providers = map[string]Provider{}
)

func Register(p Provider) {
	if p == nil {
		return
	}
	regMu.Lock()
	defer regMu.Unlock()
	providers[strings.ToLower(strings.TrimSpace(p.Name()))] = p
}

func Find(name string) Provider {
	regMu.RLock()
	defer regMu.RUnlock()
	return providers[strings.ToLower(strings.TrimSpace(name))]
}

func ResetRegistryForTest() {
	regMu.Lock()
	defer regMu.Unlock()
	providers = map[string]Provider{}
}
