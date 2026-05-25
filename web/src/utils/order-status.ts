import { i18n } from '@/locales'

const t = (key: string) => i18n.global.t(key)

export function ORDER_STATUS_LABELS(): Record<number, string> {
  return {
    1: t('orderStatus.status1'),
    2: t('orderStatus.status2'),
    3: t('orderStatus.status3'),
    4: t('orderStatus.status4'),
    5: t('orderStatus.status5'),
  }
}

export function AFTER_SALE_STATUS_LABELS(): Record<string, string> {
  return {
    applied: t('afterSaleStatus.applied'),
    approved_wait_user_return: t('afterSaleStatus.approved_wait_user_return'),
    user_returning: t('afterSaleStatus.user_returning'),
    warehouse_received: t('afterSaleStatus.warehouse_received'),
    refund_pending: t('afterSaleStatus.refund_pending'),
    refunded: t('afterSaleStatus.refunded'),
    reship_pending: t('afterSaleStatus.reship_pending'),
    reshipped: t('afterSaleStatus.reshipped'),
    completed: t('afterSaleStatus.completed'),
    rejected: t('afterSaleStatus.rejected'),
    closed: t('afterSaleStatus.closed'),
  }
}

export function SHIPMENT_STATUS_LABELS(): Record<string, string> {
  return {
    pending: t('shipmentStatus.pending'),
    shipped: t('shipmentStatus.shipped'),
    in_transit: t('shipmentStatus.in_transit'),
    signed: t('shipmentStatus.signed'),
    exception: t('shipmentStatus.exception'),
  }
}

export function SHIPMENT_BIZ_TYPE_LABELS(): Record<string, string> {
  return {
    initial: t('shipmentBizType.initial'),
    reship: t('shipmentBizType.reship'),
    return: t('shipmentBizType.return'),
  }
}

export function SHIPMENT_DIRECTION_LABELS(): Record<string, string> {
  return {
    outbound: t('shipmentDirection.outbound'),
    inbound: t('shipmentDirection.inbound'),
  }
}

export function LOGISTICS_PROVIDER_LABELS(): Record<string, string> {
  return {
    kuaidi100: t('logisticsProvider.kuaidi100'),
    kdniao: t('logisticsProvider.kdniao'),
  }
}

export function orderStatusLabel(status: number | string | undefined | null) {
  const key = Number(status || 0)
  return ORDER_STATUS_LABELS()[key] || String(status || '')
}

export function afterSaleStatusLabel(status: string | undefined | null) {
  const value = String(status || '')
  return AFTER_SALE_STATUS_LABELS()[value] || value
}

export function shipmentStatusLabel(status: string | undefined | null) {
  const value = String(status || '')
  return SHIPMENT_STATUS_LABELS()[value] || value
}

export function shipmentDirectionLabel(direction: string | undefined | null) {
  const value = String(direction || '')
  return SHIPMENT_DIRECTION_LABELS()[value] || value
}

export function shipmentBizTypeLabel(bizType: string | undefined | null) {
  const value = String(bizType || '')
  return SHIPMENT_BIZ_TYPE_LABELS()[value] || value
}

export function logisticsProviderLabel(provider: string | undefined | null) {
  const code = String(provider || '').trim()
  if (!code) return t('logisticsProvider.unbound')
  const name = LOGISTICS_PROVIDER_LABELS()[code]
  return name ? `${name}（${code}）` : code
}

export function shipmentTitle(shipment: any) {
  return `${shipmentDirectionLabel(shipment?.direction)} · ${shipmentBizTypeLabel(shipment?.biz_type)}`
}

export function shipmentPrimaryTime(shipment: any) {
  if (!shipment) return ''
  return String(shipment.signed_at || shipment.shipped_at || shipment.created_at || '')
}

export function shipmentTimeLabel(shipment: any) {
  if (shipment?.signed_at) return t('shipmentTime.signed')
  if (shipment?.shipped_at) return t('shipmentTime.shipped')
  return t('shipmentTime.recorded')
}

export function hasReshipShipment(shipments: any[]) {
  const list = Array.isArray(shipments) ? shipments : []
  return list.some((shipment: any) => String(shipment?.biz_type || '') === 'reship')
}
