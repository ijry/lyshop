import { get } from '@/utils/request'

export type DashboardData = {
  today_orders: number
  today_sales: number
  today_avg_price: number
  pending_ship: number
  pending_after_sale: number
  unread_message: number
  stock_warning: number
  compare: { revenue_yoy: number; revenue_mom: number; order_yoy: number; order_mom: number }
  trend: {
    revenue_7d: { categories: string[]; series: Array<{ name: string; data: number[] }> }
    revenue_30d: { categories: string[]; series: Array<{ name: string; data: number[] }> }
    order_7d: { categories: string[]; series: Array<{ name: string; data: number[] }> }
  }
  status_distribution: Array<{ name: string; value: number }>
  hot_products: Array<{ id: number; title: string; cover: string; sold_qty: number }>
  announcements: Array<{ id: number; title: string; content: string; type: string; created_at: string }>
  stock_warning_list: Array<{ product_id: number; sku_id: number; title: string; stock: number; threshold: number }>
}

export const getDashboard = () => get<DashboardData>('/dashboard')
