import { ref } from 'vue'
import { getDashboard, type DashboardData } from '@/api/dashboard'

const empty: DashboardData = {
  today_orders: 0, today_sales: 0, today_avg_price: 0,
  pending_ship: 0, pending_after_sale: 0, unread_message: 0, stock_warning: 0,
  compare: { revenue_yoy: 0, revenue_mom: 0, order_yoy: 0, order_mom: 0 },
  sales_trend: [],
  status_distribution: [], hot_products: [], announcements: [], stock_warning_list: [],
}

export function useDashboard() {
  const loading = ref(false)
  const data = ref<DashboardData>({ ...empty })

  async function load() {
    loading.value = true
    try {
      const ret = await getDashboard()
      if (ret) data.value = ret
    } catch {
      data.value = { ...empty }
    } finally {
      loading.value = false
    }
  }

  return { loading, data, load }
}
