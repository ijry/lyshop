import { i18n } from '@/locales'

const t = (key: string) => i18n.global.t(key)

export function orderStatusLabel(status: number | string | undefined | null) {
  const key = Number(status || 0)
  const i18nKey = `orderStatus.${key}`
  const label = t(i18nKey)
  return label !== i18nKey ? label : String(status || '')
}

export function afterSaleStatusLabel(status: string | undefined | null) {
  const value = String(status || '')
  const i18nKey = `afterSaleStatus.${value}`
  const label = t(i18nKey)
  return label !== i18nKey ? label : value
}

export function deliveryTypeLabel(type_: string | undefined | null) {
  const value = String(type_ || 'express')
  const i18nKey = `deliveryType.${value}`
  const label = t(i18nKey)
  return label !== i18nKey ? label : value
}

export function shipmentStatusLabel(status: string | undefined | null) {
  const value = String(status || '')
  const i18nKey = `shipmentStatus.${value}`
  const label = t(i18nKey)
  return label !== i18nKey ? label : value
}

export function shipmentDirectionLabel(direction: string | undefined | null) {
  const value = String(direction || '')
  const i18nKey = `shipmentDirection.${value}`
  const label = t(i18nKey)
  return label !== i18nKey ? label : value
}

export function shipmentBizTypeLabel(bizType: string | undefined | null) {
  const value = String(bizType || '')
  const i18nKey = `shipmentBizType.${value}`
  const label = t(i18nKey)
  return label !== i18nKey ? label : value
}

export function logisticsProviderLabel(provider: string | undefined | null) {
  const code = String(provider || '').trim()
  if (!code) return t('logisticsProvider.unbound')
  const i18nKey = `logisticsProvider.${code}`
  const name = t(i18nKey)
  return name !== i18nKey ? `${name}（${code}）` : code
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
  return t('shipmentTime.recorded')
}

export function hasReshipShipment(shipments: any[]) {
  const list = Array.isArray(shipments) ? shipments : []
  return list.some((shipment: any) => String(shipment?.biz_type || '') === 'reship')
}
