package service

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCompleteDraftDocOutboundSuccess(t *testing.T) {
	testDB := setupWmsTestDB(t)
	ctx := context.Background()

	warehouse := wmsmodel.Warehouse{Code: "WH-A", Name: "A仓", Status: wmsmodel.WarehouseStatusEnabled}
	require.NoError(t, testDB.Create(&warehouse).Error)
	require.NoError(t, testDB.Create(&wmsmodel.InventoryStock{
		WarehouseID: warehouse.ID,
		SkuID:       101,
		Qty:         8,
		SafeQty:     2,
	}).Error)

	doc, err := CreateDraftDoc(ctx, CreateDocInput{
		WarehouseID: warehouse.ID,
		DocType:     wmsmodel.DocTypeOutbound,
		Remark:      "测试出库",
		Items: []DocItemInput{
			{SkuID: 101, Qty: 3},
		},
	})
	require.NoError(t, err)

	require.NoError(t, CompleteDraftDoc(ctx, doc.ID))

	var stock wmsmodel.InventoryStock
	require.NoError(t, testDB.Where("warehouse_id = ? AND sku_id = ?", warehouse.ID, 101).First(&stock).Error)
	require.Equal(t, 5, stock.Qty)

	var movements []wmsmodel.InventoryMovement
	require.NoError(t, testDB.Where("doc_id = ?", doc.ID).Find(&movements).Error)
	require.Len(t, movements, 1)
	require.Equal(t, -3, movements[0].ChangeQty)
	require.Equal(t, 8, movements[0].BeforeQty)
	require.Equal(t, 5, movements[0].AfterQty)

	var completed wmsmodel.InventoryDoc
	require.NoError(t, testDB.Where("id = ?", doc.ID).First(&completed).Error)
	require.Equal(t, wmsmodel.DocStatusCompleted, completed.Status)
	require.NotNil(t, completed.CompletedAt)
}

func TestCompleteDraftDocRollbackOnInsufficientStock(t *testing.T) {
	testDB := setupWmsTestDB(t)
	ctx := context.Background()

	warehouse := wmsmodel.Warehouse{Code: "WH-B", Name: "B仓", Status: wmsmodel.WarehouseStatusEnabled}
	require.NoError(t, testDB.Create(&warehouse).Error)
	require.NoError(t, testDB.Create(&wmsmodel.InventoryStock{
		WarehouseID: warehouse.ID,
		SkuID:       201,
		Qty:         10,
	}).Error)
	require.NoError(t, testDB.Create(&wmsmodel.InventoryStock{
		WarehouseID: warehouse.ID,
		SkuID:       202,
		Qty:         1,
	}).Error)

	doc, err := CreateDraftDoc(ctx, CreateDocInput{
		WarehouseID: warehouse.ID,
		DocType:     wmsmodel.DocTypeOutbound,
		Remark:      "回滚测试",
		Items: []DocItemInput{
			{SkuID: 201, Qty: 3},
			{SkuID: 202, Qty: 5},
		},
	})
	require.NoError(t, err)

	err = CompleteDraftDoc(ctx, doc.ID)
	require.Error(t, err)
	require.ErrorContains(t, err, "库存不足")
	require.Equal(t, ErrorKindConflict, ErrorKindOf(err))

	var stock201 wmsmodel.InventoryStock
	require.NoError(t, testDB.Where("warehouse_id = ? AND sku_id = ?", warehouse.ID, 201).First(&stock201).Error)
	require.Equal(t, 10, stock201.Qty)

	var stock202 wmsmodel.InventoryStock
	require.NoError(t, testDB.Where("warehouse_id = ? AND sku_id = ?", warehouse.ID, 202).First(&stock202).Error)
	require.Equal(t, 1, stock202.Qty)

	var movements []wmsmodel.InventoryMovement
	require.NoError(t, testDB.Where("doc_id = ?", doc.ID).Find(&movements).Error)
	require.Len(t, movements, 0)

	var latestDoc wmsmodel.InventoryDoc
	require.NoError(t, testDB.Where("id = ?", doc.ID).First(&latestDoc).Error)
	require.Equal(t, wmsmodel.DocStatusDraft, latestDoc.Status)
	require.Nil(t, latestDoc.CompletedAt)
}

func TestCompleteDraftDocConflictWhenRepeated(t *testing.T) {
	testDB := setupWmsTestDB(t)
	ctx := context.Background()

	warehouse := wmsmodel.Warehouse{Code: "WH-R", Name: "重复完成仓", Status: wmsmodel.WarehouseStatusEnabled}
	require.NoError(t, testDB.Create(&warehouse).Error)
	require.NoError(t, testDB.Create(&wmsmodel.InventoryStock{
		WarehouseID: warehouse.ID,
		SkuID:       301,
		Qty:         10,
	}).Error)

	doc, err := CreateDraftDoc(ctx, CreateDocInput{
		WarehouseID: warehouse.ID,
		DocType:     wmsmodel.DocTypeOutbound,
		Items: []DocItemInput{
			{SkuID: 301, Qty: 2},
		},
	})
	require.NoError(t, err)
	require.NoError(t, CompleteDraftDoc(ctx, doc.ID))

	secondErr := CompleteDraftDoc(ctx, doc.ID)
	require.Error(t, secondErr)
	require.Equal(t, ErrorKindConflict, ErrorKindOf(secondErr))
	require.ErrorContains(t, secondErr, "单据状态非法")

	var stock wmsmodel.InventoryStock
	require.NoError(t, testDB.Where("warehouse_id = ? AND sku_id = ?", warehouse.ID, 301).First(&stock).Error)
	require.Equal(t, 8, stock.Qty)

	var movements []wmsmodel.InventoryMovement
	require.NoError(t, testDB.Where("doc_id = ?", doc.ID).Find(&movements).Error)
	require.Len(t, movements, 1)
}

func TestGenDocNoUniquenessAndFormat(t *testing.T) {
	seen := make(map[string]struct{}, 200)
	for i := 0; i < 200; i++ {
		docNo := genDocNo(wmsmodel.DocTypeOutbound)
		require.True(t, strings.HasPrefix(docNo, "WO"))
		_, exists := seen[docNo]
		require.False(t, exists, "doc_no duplicated: %s", docNo)
		seen[docNo] = struct{}{}
	}
	require.True(t, strings.HasPrefix(genDocNo(wmsmodel.DocTypeInbound), "WI"))
	require.True(t, strings.HasPrefix(genDocNo("other"), "WD"))
}

func TestCreateWarehouseDuplicateCodeConflict(t *testing.T) {
	testDB := setupWmsTestDB(t)
	ctx := context.Background()

	first := &wmsmodel.Warehouse{Code: "WH-DUP", Name: "仓库1", Status: wmsmodel.WarehouseStatusEnabled}
	require.NoError(t, CreateWarehouse(ctx, first))

	second := &wmsmodel.Warehouse{Code: "WH-DUP", Name: "仓库2", Status: wmsmodel.WarehouseStatusEnabled}
	err := CreateWarehouse(ctx, second)
	require.Error(t, err)
	require.Equal(t, ErrorKindConflict, ErrorKindOf(err))
	require.ErrorContains(t, err, "仓库编码已存在")

	var count int64
	require.NoError(t, testDB.Model(&wmsmodel.Warehouse{}).Where("code = ?", "WH-DUP").Count(&count).Error)
	require.Equal(t, int64(1), count)
}

func setupWmsTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:wms_%d?mode=memory&cache=shared", time.Now().UnixNano())
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })

	require.NoError(t, gdb.AutoMigrate(
		&wmsmodel.Warehouse{},
		&wmsmodel.InventoryStock{},
		&wmsmodel.InventoryMovement{},
		&wmsmodel.InventoryDoc{},
		&wmsmodel.InventoryDocItem{},
		&wmsmodel.InventoryReservation{},
	))
	return gdb
}
