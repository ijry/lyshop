import { i18n } from '@/locales'

const t = (key: string) => i18n.global.t(key)

export function orderStatusLabels(): Record<number, string> {
  return {
    1: t('orderStatus.pending'),
    2: t('orderStatus.shipped'),
    3: t('orderStatus.delivering'),
    4: t('orderStatus.completed'),
    5: t('orderStatus.afterSale'),
  }
}

export function afterSaleStatusLabels(): Record<string, string> {
  return {
    applied: t('afterSaleStatus.applied'),
    approved_wait_user_return: t('afterSaleStatus.waitReturn'),
    user_returning: t('afterSaleStatus.returning'),
    warehouse_received: t('afterSaleStatus.warehouseReceived'),
    refund_pending: t('afterSaleStatus.refundPending'),
    refunded: t('afterSaleStatus.refunded'),
    reship_pending: t('afterSaleStatus.reshipPending'),
    reshipped: t('afterSaleStatus.reshipped'),
    completed: t('afterSaleStatus.completed'),
    rejected: t('afterSaleStatus.rejected'),
    closed: t('afterSaleStatus.closed'),
  }
}

export function deliveryTypeLabels(): Record<string, string> {
  return {
    express: t('deliveryType.express'),
    local: t('deliveryType.local'),
  }
}

export function shipmentStatusLabels(): Record<string, string> {
  return {
    pending: t('shipmentStatus.pending'),
    shipped: t('shipmentStatus.shipped'),
    in_transit: t('shipmentStatus.inTransit'),
    signed: t('shipmentStatus.signed'),
    exception: t('shipmentStatus.exception'),
  }
}

export function shipmentBizTypeLabels(): Record<string, string> {
  return {
    initial: t('shipmentBizType.initial'),
    reship: t('shipmentBizType.reship'),
    return: t('shipmentBizType.return'),
  }
}

export function shipmentDirectionLabels(): Record<string, string> {
  return {
    outbound: t('shipmentDirection.outbound'),
    inbound: t('shipmentDirection.inbound'),
  }
}

export function orderStatusLabel(status: number | string | undefined | null) {
  const key = Number(status || 0)
  return orderStatusLabels()[key] || String(status || '')
}

export function afterSaleStatusLabel(status: string | undefined | null) {
  const value = String(status || '')
  return afterSaleStatusLabels()[value] || value
}

export function deliveryTypeLabel(type_: string | undefined | null) {
  const value = String(type_ || 'express')
  return deliveryTypeLabels()[value] || value
}

export function shipmentStatusLabel(status: string | undefined | null) {
  const value = String(status || '')
  return shipmentStatusLabels()[value] || value
}

export function shipmentDirectionLabel(direction: string | undefined | null) {
  const value = String(direction || '')
  return shipmentDirectionLabels()[value] || value
}

export function shipmentBizTypeLabel(bizType: string | undefined | null) {
  const value = String(bizType || '')
  return shipmentBizTypeLabels()[value] || value
}

export function shipmentTitle(shipment: any) {
  const dt = deliveryTypeLabel(shipment?.delivery_type)
  return `${dt} · ${shipmentDirectionLabel(shipment?.direction)} · ${shipmentBizTypeLabel(shipment?.biz_type)}`
}

export function shipmentPrimaryTime(shipment: any) {
  if (!shipment) return ''
  return String(shipment.signed_at || shipment.shipped_at || shipment.created_at || '')
}

export function shipmentTimeLabel(shipment: any) {
  if (shipment?.signed_at) return t('shipmentTime.signed')
  if (shipment?.shipped_at) return t('shipmentTime.shipped')
  return t('shipmentTime.created')
}

export function hasReshipShipment(shipments: any[]) {
  const list = Array.isArray(shipments) ? shipments : []
  return list.some((shipment: any) => String(shipment?.biz_type || '') === 'reship')
}
