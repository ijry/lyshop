import { computed, ref, unref, type Ref } from 'vue'
import {
  createActivity,
  deleteActivity,
  getActivities,
  getActivityProducts,
  type MarketingKind,
  updateActivity,
  upsertActivityProducts,
} from '@/api/marketing'
import { getProductDetail, getProducts } from '@/api/product'
import { activityState, parseSkuLabel, type ActivityPhase } from '@/utils/marketing'

/**
 * 营销活动数据层：页面只负责渲染与交互，所有请求、错误提示、冗余字段补齐都在这里。
 * kind 允许传常量或 ref，方便活动详情页从 onLoad 参数解析后再赋值。
 */
type KindSource = MarketingKind | Ref<MarketingKind>
type IdSource = number | Ref<number>

function kindOf(source: KindSource): MarketingKind {
  return unref(source) as MarketingKind
}

function idOf(source: IdSource): number {
  return Number(unref(source) || 0)
}

function messageOf(error: any): string {
  const raw = String(error?.message || error?.msg || error || '').trim()
  if (!raw) return '操作失败，请稍后重试'
  return raw.length > 20 ? `${raw.slice(0, 20)}…` : raw
}

/** 请求层已提示过的错误不重复弹 toast */
function showError(error: any) {
  if (error?.__toasted) return
  uni.showToast({ title: messageOf(error), icon: 'none' })
}

export interface ActivityRow {
  id: number
  name: string
  start_at: string
  end_at: string
  status: number
  [key: string]: any
}

/* -------------------------------------------------------------------------- */
/* 活动列表                                                                    */
/* -------------------------------------------------------------------------- */

export function useActivityCenter(kindSource: KindSource) {
  const activities = ref<ActivityRow[]>([])
  const productCounts = ref<Record<number, number>>({})
  const loading = ref(false)
  const submitting = ref(false)
  const keyword = ref('')
  const phase = ref<'all' | ActivityPhase>('all')

  const filtered = computed(() => {
    const kw = keyword.value.trim().toLowerCase()
    return activities.value.filter((item) => {
      if (phase.value !== 'all' && activityState(item).phase !== phase.value) return false
      if (kw && !String(item.name || '').toLowerCase().includes(kw)) return false
      return true
    })
  })

  const stats = computed(() => {
    const result = { total: activities.value.length, running: 0, upcoming: 0, ended: 0, disabled: 0 }
    for (const item of activities.value) result[activityState(item).phase] += 1
    return result
  })

  async function load() {
    loading.value = true
    const kind = kindOf(kindSource)
    try {
      const [activityData, productData] = await Promise.all([
        getActivities(kind, { page: 1, size: 200 }),
        getActivityProducts(kind, { page: 1, size: 200 }).catch(() => null),
      ]) as [any, any]
      activities.value = Array.isArray(activityData?.list)
        ? activityData.list
        : Array.isArray(activityData)
          ? activityData
          : []
      const productList = Array.isArray(productData?.list)
        ? productData.list
        : Array.isArray(productData)
          ? productData
          : []
      const counts: Record<number, number> = {}
      for (const row of productList) {
        const id = Number(row?.activity_id || 0)
        if (id) counts[id] = (counts[id] || 0) + 1
      }
      productCounts.value = counts
    } catch (error) {
      activities.value = []
      productCounts.value = {}
      showError(error)
    } finally {
      loading.value = false
    }
  }

  async function save(payload: any, id = 0): Promise<boolean> {
    if (submitting.value) return false
    submitting.value = true
    try {
      const kind = kindOf(kindSource)
      if (id > 0) await updateActivity(kind, id, payload)
      else await createActivity(kind, payload)
      await load()
      uni.showToast({ title: id > 0 ? '活动已更新' : '活动已创建', icon: 'success' })
      return true
    } catch (error) {
      showError(error)
      return false
    } finally {
      submitting.value = false
    }
  }

  async function toggle(item: ActivityRow): Promise<boolean> {
    const id = Number(item?.id || 0)
    if (!id) return false
    return save(
      {
        name: String(item.name || ''),
        start_at: String(item.start_at || ''),
        end_at: String(item.end_at || ''),
        status: Number(item.status || 0) === 1 ? 0 : 1,
      },
      id,
    )
  }

  async function remove(id: number): Promise<boolean> {
    if (submitting.value) return false
    submitting.value = true
    try {
      await deleteActivity(kindOf(kindSource), id)
      await load()
      uni.showToast({ title: '活动已删除', icon: 'success' })
      return true
    } catch (error) {
      showError(error)
      return false
    } finally {
      submitting.value = false
    }
  }

  return {
    activities,
    productCounts,
    loading,
    submitting,
    keyword,
    phase,
    filtered,
    stats,
    load,
    save,
    toggle,
    remove,
  }
}

/* -------------------------------------------------------------------------- */
/* 活动商品                                                                    */
/* -------------------------------------------------------------------------- */

export interface ActivityProductRow {
  id?: number
  product_id: number
  sku_id: number
  activity_price: number
  start_price: number
  floor_price: number
  limit_per_order: number
  total_stock_limit: number
  sold_qty?: number
  product_title?: string
  product_cover?: string
  sku_price?: number
  sku_stock?: number
  sku_attrs?: any
  [key: string]: any
}

export interface SkuOption {
  id: number
  label: string
  price: number
  stock: number
  attrs: any
}

/** 后端 PUT 会整体替换活动商品，必须把展示所需的冗余字段一并回传 */
function toPayloadRow(row: ActivityProductRow) {
  return {
    product_id: Number(row.product_id || 0),
    sku_id: Number(row.sku_id || 0),
    activity_price: Number(row.activity_price || 0),
    start_price: Number(row.start_price || 0),
    floor_price: Number(row.floor_price || 0),
    limit_per_order: Number(row.limit_per_order || 0),
    total_stock_limit: Number(row.total_stock_limit || 0),
    sold_qty: Number(row.sold_qty || 0),
    product_title: String(row.product_title || ''),
    product_cover: String(row.product_cover || ''),
    sku_price: Number(row.sku_price || 0),
    sku_stock: Number(row.sku_stock || 0),
    sku_attrs: Array.isArray(row.sku_attrs) ? row.sku_attrs : [],
  }
}

export function useActivityProducts(kindSource: KindSource, activityIDSource: IdSource) {
  const rows = ref<ActivityProductRow[]>([])
  const loading = ref(false)
  const submitting = ref(false)
  const keyword = ref('')
  const productOptions = ref<any[]>([])
  const skuOptions = ref<SkuOption[]>([])

  const filtered = computed(() => {
    const kw = keyword.value.trim().toLowerCase()
    if (!kw) return rows.value
    return rows.value.filter((item) => String(item.product_title || '').toLowerCase().includes(kw))
  })

  const stats = computed(() => {
    let soldQty = 0
    let gross = 0
    for (const row of rows.value) {
      const qty = Number(row.sold_qty || 0)
      soldQty += qty
      gross += Number(row.activity_price || row.start_price || 0) * qty
    }
    return { count: rows.value.length, soldQty, gross }
  })

  async function load() {
    const id = idOf(activityIDSource)
    if (!id) {
      rows.value = []
      return
    }
    loading.value = true
    try {
      const data: any = await getActivityProducts(kindOf(kindSource), {
        activity_id: id,
        page: 1,
        size: 200,
      })
      rows.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
    } catch (error) {
      rows.value = []
      showError(error)
    } finally {
      loading.value = false
    }
  }

  async function loadProducts() {
    try {
      const data: any = await getProducts({ page: 1, size: 200 })
      productOptions.value = Array.isArray(data?.list) ? data.list : Array.isArray(data) ? data : []
    } catch (error) {
      productOptions.value = []
      showError(error)
    }
  }

  async function loadSkus(productID: number): Promise<SkuOption[]> {
    skuOptions.value = []
    if (!productID) return []
    try {
      const detail: any = await getProductDetail(productID)
      const skus = Array.isArray(detail?.skus) ? detail.skus : []
      skuOptions.value = skus
        .map((item: any) => ({
          id: Number(item?.id || 0),
          label: `${parseSkuLabel(item?.attrs, Number(item?.id || 0))} ｜ ¥${Number(item?.price || 0).toFixed(2)}`,
          price: Number(item?.price || 0),
          stock: Number(item?.stock || 0),
          attrs: item?.attrs || [],
        }))
        .filter((item) => item.id > 0)
    } catch (error) {
      showError(error)
    }
    return skuOptions.value
  }

  /** 新建行时补齐封面 / 标题 / 原价，避免整表替换后卡片退化成「商品#ID」 */
  function enrich(row: ActivityProductRow): Partial<ActivityProductRow> {
    const product = productOptions.value.find((item) => Number(item.id || 0) === Number(row.product_id || 0))
    const sku = skuOptions.value.find((item) => Number(item.id || 0) === Number(row.sku_id || 0))
    return {
      product_title: String(product?.title || row.product_title || `商品 #${row.product_id}`),
      product_cover: String(product?.cover || row.product_cover || ''),
      sku_price: Number(sku?.price ?? row.sku_price ?? 0),
      sku_stock: Number(sku?.stock ?? row.sku_stock ?? 0),
      sku_attrs: sku?.attrs ?? row.sku_attrs ?? [],
    }
  }

  async function save(row: ActivityProductRow, editingID = 0): Promise<boolean> {
    const activityID = idOf(activityIDSource)
    if (!activityID || submitting.value) return false
    const duplicate = rows.value.find(
      (item) =>
        Number(item.id || 0) !== editingID &&
        Number(item.product_id || 0) === Number(row.product_id || 0) &&
        Number(item.sku_id || 0) === Number(row.sku_id || 0),
    )
    if (duplicate) {
      uni.showToast({ title: '同一活动下 SKU 不可重复', icon: 'none' })
      return false
    }
    submitting.value = true
    try {
      const previous = rows.value.find((item) => Number(item.id || 0) === editingID)
      const merged: ActivityProductRow = {
        ...(previous || {}),
        ...row,
        ...enrich(row),
      }
      const nextRows = rows.value.filter((item) => Number(item.id || 0) !== editingID)
      nextRows.push(merged)
      await upsertActivityProducts(kindOf(kindSource), activityID, nextRows.map(toPayloadRow))
      await load()
      uni.showToast({ title: editingID ? '商品已更新' : '商品已添加', icon: 'success' })
      return true
    } catch (error) {
      showError(error)
      return false
    } finally {
      submitting.value = false
    }
  }

  async function remove(rowID: number): Promise<boolean> {
    const activityID = idOf(activityIDSource)
    if (!activityID || submitting.value) return false
    submitting.value = true
    try {
      const nextRows = rows.value.filter((item) => Number(item.id || 0) !== Number(rowID || 0))
      await upsertActivityProducts(kindOf(kindSource), activityID, nextRows.map(toPayloadRow))
      await load()
      uni.showToast({ title: '商品已移除', icon: 'success' })
      return true
    } catch (error) {
      showError(error)
      return false
    } finally {
      submitting.value = false
    }
  }

  return {
    rows,
    loading,
    submitting,
    keyword,
    productOptions,
    skuOptions,
    filtered,
    stats,
    load,
    loadProducts,
    loadSkus,
    save,
    remove,
  }
}
