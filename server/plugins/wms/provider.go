package wms

import (
	"context"

	inventorycore "github.com/ijry/lyshop/core/inventory"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	wmssvc "github.com/ijry/lyshop/plugins/wms/service"
	"gorm.io/gorm"
)

type wmssvcReservationItemInput = wmssvc.ReservationItemInput
type wmssvcReserveStockInput = wmssvc.ReserveStockInput

var (
	reserveStockTxFn           = wmssvc.ReserveStockTx
	confirmReservationTxFn     = wmssvc.ConfirmReservationTx
	releaseReservationTxFn     = wmssvc.ReleaseReservationTx
	pickDefaultWarehouseIDTxFn = wmssvc.PickDefaultWarehouseIDTx
	listStocksBySkuIDsFn       = wmssvc.ListStocksBySkuIDs
	createDraftDocFn           = wmssvc.CreateDraftDoc
	completeDraftDocFn         = wmssvc.CompleteDraftDoc
)

type builtinProvider struct{}

func (p *builtinProvider) Name() string { return "builtin_wms" }

func (p *builtinProvider) ReserveTx(tx *gorm.DB, in inventorycore.ReserveInput) error {
	items := make([]wmssvcReservationItemInput, 0, len(in.Items))
	for _, item := range in.Items {
		items = append(items, wmssvcReservationItemInput{SkuID: item.SkuID, Qty: item.Qty})
	}
	warehouseID, err := pickDefaultWarehouseIDTxFn(tx)
	if err != nil {
		return err
	}
	return reserveStockTxFn(tx, wmssvcReserveStockInput{
		BizType:     in.BizType,
		BizNo:       in.BizNo,
		WarehouseID: warehouseID,
		Items:       items,
		ExpiredAt:   in.ExpiredAt,
	})
}

func (p *builtinProvider) ConfirmTx(tx *gorm.DB, bizType, bizNo string) error {
	return confirmReservationTxFn(tx, bizType, bizNo)
}

func (p *builtinProvider) ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error {
	return releaseReservationTxFn(tx, bizType, bizNo, reason)
}

func (p *builtinProvider) DeductTx(tx *gorm.DB, in inventorycore.DeductInput) error {
	warehouseID, err := pickDefaultWarehouseIDTxFn(tx)
	if err != nil {
		return err
	}
	doc := wmssvc.CreateDocInput{
		WarehouseID: warehouseID,
		DocType:     wmsmodel.DocTypeOutbound,
		Remark:      in.BizNo,
		Items:       make([]wmssvc.DocItemInput, 0, len(in.Items)),
	}
	for _, item := range in.Items {
		doc.Items = append(doc.Items, wmssvc.DocItemInput{SkuID: item.SkuID, Qty: item.Qty})
	}
	row, err := createDraftDocFn(context.Background(), doc)
	if err != nil {
		return err
	}
	return completeDraftDocFn(context.Background(), row.ID)
}

func (p *builtinProvider) RestoreTx(tx *gorm.DB, in inventorycore.RestoreInput) error {
	warehouseID, err := pickDefaultWarehouseIDTxFn(tx)
	if err != nil {
		return err
	}
	doc := wmssvc.CreateDocInput{
		WarehouseID: warehouseID,
		DocType:     wmsmodel.DocTypeInbound,
		Remark:      in.BizNo,
		Items:       make([]wmssvc.DocItemInput, 0, len(in.Items)),
	}
	for _, item := range in.Items {
		doc.Items = append(doc.Items, wmssvc.DocItemInput{SkuID: item.SkuID, Qty: item.Qty})
	}
	row, err := createDraftDocFn(context.Background(), doc)
	if err != nil {
		return err
	}
	return completeDraftDocFn(context.Background(), row.ID)
}

func (p *builtinProvider) SyncSkuTx(tx *gorm.DB, in inventorycore.SyncSkuInput) error {
	warehouseID, err := pickDefaultWarehouseIDTxFn(tx)
	if err != nil || warehouseID == 0 {
		return err
	}
	var stock wmsmodel.InventoryStock
	result := tx.Where("warehouse_id = ? AND sku_id = ?", warehouseID, in.SkuID).First(&stock)
	if result.Error == nil {
		return nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}
	return tx.Create(&wmsmodel.InventoryStock{
		WarehouseID: warehouseID,
		SkuID:       in.SkuID,
		Qty:         in.Stock,
	}).Error
}

func (p *builtinProvider) GetSellableStock(ctx context.Context, skuIDs []uint64) ([]inventorycore.SellableStock, error) {
	rows, err := listStocksBySkuIDsFn(ctx, skuIDs)
	if err != nil {
		return nil, err
	}
	out := make([]inventorycore.SellableStock, 0, len(rows))
	for _, row := range rows {
		out = append(out, inventorycore.SellableStock{
			SkuID:    row.SkuID,
			OnHand:   row.Qty,
			Reserved: row.ReservedQty,
			Sellable: row.Qty - row.ReservedQty,
		})
	}
	return out, nil
}
