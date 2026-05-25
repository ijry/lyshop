export const ORDER_STATUS_LABELS: Record<number, string> = {
  1: '待付款',
  2: '待发货',
  3: '待收货',
  4: '已完成',
  5: '售后',
}

export const AFTER_SALE_STATUS_LABELS: Record<string, string> = {
  applied: '已申请',
  approved_wait_user_return: '待用户回寄',
  user_returning: '用户回寄中',
  warehouse_received: '仓库已收货',
  refund_pending: '待退款',
  refunded: '已退款',
  reship_pending: '待补发',
  reshipped: '已补发',
  completed: '已完结',
  rejected: '已拒绝',
  closed: '已关闭',
}

export const DELIVERY_TYPE_LABELS: Record<string, string> = {
  express: '快递配送',
  local: '同城配送',
}

export const SHIPMENT_STATUS_LABELS: Record<string, string> = {
  pending: '待揽收',
  shipped: '已发货',
  in_transit: '运输中',
  signed: '已签收',
  exception: '物流异常',
}

export const SHIPMENT_BIZ_TYPE_LABELS: Record<string, string> = {
  initial: '首发',
  reship: '补发',
  return: '回寄',
}

export const SHIPMENT_DIRECTION_LABELS: Record<string, string> = {
  outbound: '寄出',
  inbound: '回寄',
}

export function orderStatusLabel(status: number | string | undefined | null) {
  const key = Number(status || 0)
  return ORDER_STATUS_LABELS[key] || String(status || '')
}

export function afterSaleStatusLabel(status: string | undefined | null) {
  const value = String(status || '')
  return AFTER_SALE_STATUS_LABELS[value] || value
}

export function deliveryTypeLabel(type_: string | undefined | null) {
  const value = String(type_ || 'express')
  return DELIVERY_TYPE_LABELS[value] || value
}

export function shipmentStatusLabel(status: string | undefined | null) {
  const value = String(status || '')
  return SHIPMENT_STATUS_LABELS[value] || value
}

export function shipmentDirectionLabel(direction: string | undefined | null) {
  const value = String(direction || '')
  return SHIPMENT_DIRECTION_LABELS[value] || value
}

export function shipmentBizTypeLabel(bizType: string | undefined | null) {
  const value = String(bizType || '')
  return SHIPMENT_BIZ_TYPE_LABELS[value] || value
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
  if (shipment?.signed_at) return '签收时间'
  if (shipment?.shipped_at) return '发货时间'
  return '记录时间'
}

export function hasReshipShipment(shipments: any[]) {
  const list = Array.isArray(shipments) ? shipments : []
  return list.some((shipment: any) => String(shipment?.biz_type || '') === 'reship')
}
