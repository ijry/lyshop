package service

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	decormodel "github.com/ijry/lyshop/plugins/decor/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDefaultComponentsForPagePC(t *testing.T) {
	got := defaultComponentsForPage("pc")
	require.JSONEq(t, `{"pageStyle":{"background":{"mode":"solid","solidColor":"#f8fafc","overlay":{"enabled":false,"color":"#000000","opacity":0.2}},"content":{"maxWidth":1280,"gutterX":24,"sectionGap":24},"surface":{"radius":12,"shadow":"none"}},"components":[]}`, string(got))
}

func TestDefaultComponentsForPageNonPC(t *testing.T) {
	got := defaultComponentsForPage("index")
	require.JSONEq(t, `[]`, string(got))
}

func setupDecorServiceTestDB(t *testing.T) {
	t.Helper()
	dbName := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	gdb, err := gorm.Open(sqlite.Open("file:"+dbName+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, gdb.AutoMigrate(&decormodel.DecorPage{}))
	db.DB = gdb
}

func TestListVariantsForPCReturnsDefaultPayload(t *testing.T) {
	setupDecorServiceTestDB(t)

	rows, err := ListVariants(context.Background(), 0, "pc")

	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "pc", rows[0].PageKey)
	require.Equal(t, DefaultVariantKey, rows[0].VariantKey)
	require.JSONEq(t, `{"pageStyle":{"background":{"mode":"solid","solidColor":"#f8fafc","overlay":{"enabled":false,"color":"#000000","opacity":0.2}},"content":{"maxWidth":1280,"gutterX":24,"sectionGap":24},"surface":{"radius":12,"shadow":"none"}},"components":[]}`, string(rows[0].Components))
}

func TestCreateVariantCopyForPCKeepsPagePayload(t *testing.T) {
	setupDecorServiceTestDB(t)
	ctx := context.Background()
	payload := json.RawMessage(`{"pageStyle":{"background":{"mode":"solid","solidColor":"#111111"},"content":{"maxWidth":1180,"gutterX":20,"sectionGap":18},"surface":{"radius":8,"shadow":"sm"}},"components":[{"id":"pc_hero","type":"hero","props":{"title":"A"}}]}`)
	_, err := SavePage(ctx, 0, "pc", payload, "default")
	require.NoError(t, err)

	row, err := CreateVariantCopy(ctx, 0, "pc", "default", "summer", "夏季版")

	require.NoError(t, err)
	require.Equal(t, "pc", row.PageKey)
	require.Equal(t, "summer", row.VariantKey)
	require.Equal(t, "夏季版", row.VariantName)
	require.False(t, row.IsPublished)
	require.JSONEq(t, string(payload), string(row.Components))
}

func TestPublishPageForPCIsSingleActiveWithinPCOnly(t *testing.T) {
	setupDecorServiceTestDB(t)
	ctx := context.Background()
	pcPayload := json.RawMessage(`{"pageStyle":{"background":{"mode":"solid","solidColor":"#f8fafc"},"content":{"maxWidth":1280,"gutterX":24,"sectionGap":24},"surface":{"radius":12,"shadow":"none"}},"components":[]}`)
	indexPayload := json.RawMessage(`[{"id":"m_1","type":"banner","props":{}}]`)
	_, err := SavePage(ctx, 0, "pc", pcPayload, "default")
	require.NoError(t, err)
	_, err = CreateVariantCopy(ctx, 0, "pc", "default", "festival", "节日版")
	require.NoError(t, err)
	_, err = SavePage(ctx, 0, "index", indexPayload, "default")
	require.NoError(t, err)
	require.NoError(t, PublishPage(ctx, 0, "index", "default"))
	require.NoError(t, PublishPage(ctx, 0, "pc", "default"))

	require.NoError(t, PublishPage(ctx, 0, "pc", "festival"))

	var pcRows []decormodel.DecorPage
	require.NoError(t, db.DB.Where("page_key = ?", "pc").Order("variant_key asc").Find(&pcRows).Error)
	require.Len(t, pcRows, 2)
	require.False(t, pcRows[0].IsPublished)
	require.Equal(t, "default", pcRows[0].VariantKey)
	require.True(t, pcRows[1].IsPublished)
	require.Equal(t, "festival", pcRows[1].VariantKey)

	indexPublished, err := GetPublishedPage(ctx, 0, "index")
	require.NoError(t, err)
	require.Equal(t, "index", indexPublished.PageKey)
	require.Equal(t, "default", indexPublished.VariantKey)
	require.True(t, indexPublished.IsPublished)
}

