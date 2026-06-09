package inventory

import (
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestLocalProviderReserveConfirmRelease(t *testing.T) {
	testDB := setupLocalInventoryDB(t)
	p := &localProvider{}

	require.NoError(t, testDB.Create(&productmodel.ProductSku{
		ProductID: 1,
		SkuKey:    "red:l",
		Price:     10,
		Stock:     20,
		Status:    productmodel.ProductSkuStatusActive,
	}).Error)

	var sku productmodel.ProductSku
	require.NoError(t, testDB.First(&sku).Error)

	expireAt := time.Now().Add(15 * time.Minute)
	err := testDB.Transaction(func(tx *gorm.DB) error {
		return p.ReserveTx(tx, ReserveInput{
			BizType:   "order",
			BizNo:     "ORD-1",
			Items:     []ReserveItem{{SkuID: sku.ID, Qty: 3}},
			ExpiredAt: &expireAt,
		})
	})
	require.NoError(t, err)

	err = testDB.Transaction(func(tx *gorm.DB) error {
		return p.ConfirmTx(tx, "order", "ORD-1")
	})
	require.NoError(t, err)

	require.NoError(t, testDB.First(&sku, sku.ID).Error)
	require.Equal(t, 17, sku.Stock)

	err = testDB.Transaction(func(tx *gorm.DB) error {
		return p.ReleaseTx(tx, "order", "ORD-1", "already_confirmed")
	})
	require.ErrorContains(t, err, "confirmed")
}

func TestLocalProviderReserveRejectsInsufficientStock(t *testing.T) {
	testDB := setupLocalInventoryDB(t)
	p := &localProvider{}
	require.NoError(t, testDB.Create(&productmodel.ProductSku{
		ProductID: 2,
		SkuKey:    "blue:m",
		Price:     10,
		Stock:     1,
		Status:    productmodel.ProductSkuStatusActive,
	}).Error)
	var sku productmodel.ProductSku
	require.NoError(t, testDB.First(&sku).Error)
	err := testDB.Transaction(func(tx *gorm.DB) error {
		return p.ReserveTx(tx, ReserveInput{
			BizType: "order",
			BizNo:   "ORD-2",
			Items:   []ReserveItem{{SkuID: sku.ID, Qty: 2}},
		})
	})
	require.ErrorContains(t, err, "库存不足")
}

func setupLocalInventoryDB(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:local_inventory_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)
	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })
	require.NoError(t, gdb.AutoMigrate(&productmodel.ProductSku{}, &InventoryReservation{}, &OrderInventoryState{}))
	return gdb
}
