package service

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCanonicalSkuKeySorted(t *testing.T) {
	key := CanonicalSkuKey([]productmodel.SkuAttr{
		{Name: "尺码", Value: "M"},
		{Name: "颜色", Value: "红"},
	})
	require.Equal(t, "尺码:M|颜色:红", key)
}

func TestBuildSkusFromSpecSchema(t *testing.T) {
	price := 19.9
	stock := 8
	rows, err := BuildSkusFromSpecSchema(
		[]SpecSchemaGroup{
			{Name: "颜色", Values: []string{"红", "蓝"}},
			{Name: "尺码", Values: []string{"M"}},
		},
		99,
		[]SkuOverride{
			{
				SkuKey:  "尺码:M|颜色:红",
				Price:   &price,
				Stock:   &stock,
				SkuCode: "SKU-RED-M",
			},
		},
	)
	require.NoError(t, err)
	require.Len(t, rows, 2)

	byKey := map[string]productmodel.ProductSku{}
	for _, row := range rows {
		byKey[row.SkuKey] = row
	}
	require.Equal(t, "SKU-RED-M", byKey["尺码:M|颜色:红"].SkuCode)
	require.Equal(t, 19.9, byKey["尺码:M|颜色:红"].Price)
	require.Equal(t, 8, byKey["尺码:M|颜色:红"].Stock)
	require.Equal(t, 99.0, byKey["尺码:M|颜色:蓝"].Price)
}

func TestReplaceProductSkusInactivateMissing(t *testing.T) {
	testDB := setupProductTestDB(t)
	ctx := t.Context()

	prod := &productmodel.Product{
		Title:  "测试商品",
		Price:  50,
		Status: 1,
		Detail: json.RawMessage(`{"version":1,"blocks":[]}`),
	}
	require.NoError(t, CreateProduct(ctx, prod, []productmodel.ProductSku{
		mustSKU(t, []productmodel.SkuAttr{{Name: "颜色", Value: "红"}}, 50, 10),
		mustSKU(t, []productmodel.SkuAttr{{Name: "颜色", Value: "蓝"}}, 50, 10),
	}, nil))

	diff, err := ReplaceProductSkus(ctx, prod.ID, []productmodel.ProductSku{
		mustSKU(t, []productmodel.SkuAttr{{Name: "颜色", Value: "红"}}, 60, 5),
	})
	require.NoError(t, err)
	require.Equal(t, 0, diff.Added)
	require.Equal(t, 1, diff.Kept)
	require.Equal(t, 1, diff.Inactivated)

	var activeCount int64
	require.NoError(t, testDB.Model(&productmodel.ProductSku{}).
		Where("product_id = ? AND status = ?", prod.ID, productmodel.ProductSkuStatusActive).
		Count(&activeCount).Error)
	require.Equal(t, int64(1), activeCount)

	var inactiveCount int64
	require.NoError(t, testDB.Model(&productmodel.ProductSku{}).
		Where("product_id = ? AND status = ?", prod.ID, productmodel.ProductSkuStatusInactive).
		Count(&inactiveCount).Error)
	require.Equal(t, int64(1), inactiveCount)
}

func TestProductSkuKeyUniquePerProduct(t *testing.T) {
	testDB := setupProductTestDB(t)
	first := &productmodel.Product{
		Title:  "商品1",
		Price:  10,
		Status: 1,
		Detail: json.RawMessage(`{"version":1,"blocks":[]}`),
	}
	second := &productmodel.Product{
		Title:  "商品2",
		Price:  20,
		Status: 1,
		Detail: json.RawMessage(`{"version":1,"blocks":[]}`),
	}
	require.NoError(t, testDB.Create(first).Error)
	require.NoError(t, testDB.Create(second).Error)

	redAttrs, err := EncodeSkuAttrs([]productmodel.SkuAttr{{Name: "颜色", Value: "红"}})
	require.NoError(t, err)
	key := CanonicalSkuKey([]productmodel.SkuAttr{{Name: "颜色", Value: "红"}})

	require.NoError(t, testDB.Create(&productmodel.ProductSku{
		ProductID: first.ID,
		Attrs:     redAttrs,
		Price:     10,
		Stock:     5,
		SkuKey:    key,
		Status:    productmodel.ProductSkuStatusActive,
	}).Error)
	require.NoError(t, testDB.Create(&productmodel.ProductSku{
		ProductID: second.ID,
		Attrs:     redAttrs,
		Price:     20,
		Stock:     8,
		SkuKey:    key,
		Status:    productmodel.ProductSkuStatusActive,
	}).Error)

	require.Error(t, testDB.Create(&productmodel.ProductSku{
		ProductID: first.ID,
		Attrs:     redAttrs,
		Price:     10,
		Stock:     1,
		SkuKey:    key,
		Status:    productmodel.ProductSkuStatusActive,
	}).Error)
}

func mustSKU(t *testing.T, attrs []productmodel.SkuAttr, price float64, stock int) productmodel.ProductSku {
	t.Helper()
	raw, err := json.Marshal(attrs)
	require.NoError(t, err)
	return productmodel.ProductSku{
		Attrs: raw,
		Price: price,
		Stock: stock,
	}
}

func setupProductTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:product_%d?mode=memory&cache=shared", time.Now().UnixNano())
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })

	require.NoError(t, gdb.AutoMigrate(
		&productmodel.Product{},
		&productmodel.ProductSku{},
	))
	return gdb
}
