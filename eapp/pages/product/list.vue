<script setup lang="ts">
import { onLoad, onPullDownRefresh, onReachBottom, onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import PageHeader from '@/components/biz/PageHeader.vue'
import ProductCard from '@/components/biz/ProductCard.vue'
import EmptyState from '@/components/biz/EmptyState.vue'
import FilterDrawer from '@/components/biz/FilterDrawer.vue'
import BatchBar from '@/components/biz/BatchBar.vue'
import BatchResultPopup from '@/components/biz/BatchResultPopup.vue'
import AppIcon from '@/components/common/AppIcon.vue'
import { useProductList } from '@/composables/useProductList'
import { batchUpdateProductCategory, batchUpdateProductPrice, batchUpdateProductStatus, updateProduct } from '@/api/product'
import { getStorage, removeStorage } from '@/utils/storage'

const h = useProductList()
const tabs = [
  { name: '全部', status: '' as ''|0|1, low_stock: false },
  { name: '在售', status: 1 as ''|0|1, low_stock: false },
  { name: '仓库', status: 0 as ''|0|1, low_stock: false },
  { name: '预警', status: '' as ''|0|1, low_stock: true },
]
const current = ref(0)
const quickKeyword = ref('')
const showFilter = ref(false); const showResult = ref(false)
const showSortPicker = ref(false)
const sortOptions = [
  { label: '默认', value: '' },
  { label: '销量', value: 'sales' },
  { label: '库存', value: 'stock' },
  { label: '价格升', value: 'price_asc' },
  { label: '价格降', value: 'price_desc' },
  { label: '最新', value: 'created' },
]
const filterDraft = ref<any>({ keyword: '', sort_by: '' })
const result = ref<{ success_ids: number[]; fail: Array<{ id: number; reason: string }> }>({ success_ids: [], fail: [] })

function onTab(idx: number) {
  current.value = idx
  h.applyFilter({ status: tabs[idx].status, low_stock: tabs[idx].low_stock, keyword: quickKeyword.value || undefined })
  h.load()
}

function onQuickSearch() {
  h.applyFilter({ keyword: quickKeyword.value || undefined })
  h.load()
}

function clearQuickSearch() {
  quickKeyword.value = ''
  h.applyFilter({ keyword: undefined })
  h.load()
}

function openFilter() { filterDraft.value = { keyword: h.filter.value.keyword || '', sort_by: h.filter.value.sort_by || '' }; showFilter.value = true }
function applyFilter() {
  h.applyFilter({ keyword: filterDraft.value.keyword || undefined, sort_by: filterDraft.value.sort_by || undefined })
  showFilter.value = false; h.load()
}
function resetFilter() { filterDraft.value = { keyword: '', sort_by: '' } }

async function onCardAction(p: any, key: string) {
  if (key === 'edit') uni.navigateTo({ url: `/pages/product/edit?id=${p.id}` })
  else if (key === 'toggle_sale') {
    const next = Number(p.status || 0) === 1 ? 0 : 1
    await updateProduct(p.id, { product: { status: next } })
    uni.showToast({ title: '状态已更新', icon: 'success' })
    await h.refresh()
  }
}

async function startBatch(action: 'shelf_on'|'shelf_off'|'category'|'price') {
  if (!h.selectCount.value) { uni.showToast({ title: '请先勾选商品', icon: 'none' }); return }
  if (action === 'shelf_on' || action === 'shelf_off') {
    const r: any = await batchUpdateProductStatus({ ids: h.selectedIds.value, status: action === 'shelf_on' ? 1 : 0 })
    result.value = r; showResult.value = true; h.clearSelect(); await h.refresh()
  } else if (action === 'category') {
    uni.showModal({ title: '批量分类', editable: true, placeholderText: '输入新分类 ID', success: async (m) => {
      if (!m.confirm || !m.content) return
      const r: any = await batchUpdateProductCategory({ ids: h.selectedIds.value, category_id: Number(m.content) })
      result.value = r; showResult.value = true; h.clearSelect(); await h.refresh()
    } })
  } else if (action === 'price') {
    uni.showModal({ title: '批量调价（百分比）', editable: true, placeholderText: '例如 -10 表示降 10%', success: async (m) => {
      if (!m.confirm || !m.content) return
      const r: any = await batchUpdateProductPrice({ ids: h.selectedIds.value, adjustment: { type: 'percent', value: Number(m.content) } })
      result.value = r; showResult.value = true; h.clearSelect(); await h.refresh()
    } })
  }
}

function newProduct() { uni.navigateTo({ url: '/pages/product/edit?id=0' }) }
function aiHint() { uni.showModal({ title: '提示', content: 'AI 生图请到管理后台 AI 工作流操作', showCancel: false }) }
function goCategoryTree() { uni.navigateTo({ url: '/pages/product/category-tree' }) }

onLoad(() => {
  const w = String(getStorage('eapp_product_status_filter') || '')
  if (w === 'warning') { current.value = 3; removeStorage('eapp_product_status_filter') }
  onTab(current.value)
})
onShow(() => h.load())
onPullDownRefresh(async () => { await h.refresh(); uni.stopPullDownRefresh() })
onReachBottom(() => h.loadMore())
</script>

<template>
  <view class="page">
    <PageHeader title="商品管理" eyebrow="PRODUCTS">
      <template #right>
        <view class="head-action" hover-class="eapp-hover" @click="goCategoryTree">
          <AppIcon name="layers" :size="26" color="primary" />
          <text class="head-action__text">分类</text>
        </view>
        <view class="head-action" hover-class="eapp-hover" @click="openFilter">
          <AppIcon name="filter" :size="26" color="primary" />
          <text class="head-action__text">筛选</text>
        </view>
      </template>
    </PageHeader>

    <view class="tab-wrap">
      <up-tabs
        :list="tabs"
        :current="current"
        :scrollable="true"
        keyName="name"
        @click="(item) => onTab(item.index)"
        :activeStyle="{ color: 'var(--eapp-primary)', fontWeight: '700', fontSize: '28rpx' }"
        :inactiveStyle="{ color: 'var(--eapp-text-muted)', fontSize: '26rpx' }"
        :itemStyle="{ padding: '0 8rpx', height: '80rpx' }"
        lineColor="var(--eapp-primary)"
        lineWidth="40rpx"
        lineHeight="4rpx"
        :lineStyle="{ borderRadius: '2rpx' }"
      />
    </view>

    <view class="search-bar">
      <up-search
        v-model="quickKeyword"
        placeholder="搜索商品名称"
        :showAction="!!quickKeyword"
        actionText="清空"
        shape="round"
        :inputStyle="{ fontSize: '26rpx' }"
        @search="onQuickSearch"
        @custom="clearQuickSearch"
      />
    </view>

    <view class="list">
      <EmptyState v-if="!h.loading.value && !h.list.value.length" title="暂无商品" desc="调整筛选条件或新建商品" icon="product">
        <up-button size="mini" plain hairline @click="newProduct">新建商品</up-button>
      </EmptyState>
      <ProductCard
        v-for="p in h.list.value"
        :key="p.id"
        :product="p"
        :selectable="true"
        :selected="h.isSelected(p.id)"
        @toggle="h.toggleSelect(p.id)"
        @click="uni.navigateTo({ url: `/pages/product/edit?id=${p.id}` })"
        @action="(key) => onCardAction(p, key)"
      />
    </view>

    <BatchBar
      :count="h.selectCount.value"
      :actions="[
        { key: 'shelf_on', label: '批量上架', tone: 'primary' },
        { key: 'shelf_off', label: '批量下架', tone: 'warning' },
        { key: 'category', label: '批量分类', tone: 'primary' },
        { key: 'price', label: '批量调价', tone: 'primary' },
      ]"
      @action="startBatch"
      @cancel="h.clearSelect()"
    />

    <FilterDrawer :show="showFilter" title="商品筛选" @close="showFilter = false" @reset="resetFilter" @confirm="applyFilter">
      <up-input v-model="filterDraft.keyword" placeholder="商品名" />
      <view class="picker mt" @click="showSortPicker = true">排序：{{ { '': '默认', sales: '销量', stock: '库存', price_asc: '价格升', price_desc: '价格降', created: '最新' }[filterDraft.sort_by || ''] }}</view>
      <up-picker :show="showSortPicker" :columns="[sortOptions]" keyName="label" @confirm="(e) => { filterDraft.sort_by = e.value[0].value; showSortPicker = false }" @cancel="showSortPicker = false" @close="showSortPicker = false" />
    </FilterDrawer>

    <BatchResultPopup :show="showResult" :success="result.success_ids" :fails="result.fail" @close="showResult = false" />

    <view class="fab">
      <up-button shape="circle" plain @click="aiHint">AI 生图</up-button>
      <up-button type="primary" shape="circle" @click="newProduct">+ 新建商品</up-button>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); }

.head-action {
  display: flex; align-items: center; gap: 6rpx;
  padding: 10rpx 14rpx;
  border: 1rpx solid var(--eapp-primary-border);
  border-radius: var(--eapp-radius);
  background: var(--eapp-primary-soft);
  box-sizing: border-box;
}
.head-action__text { font-size: 24rpx; font-weight: 600; color: var(--eapp-primary); }

.tab-wrap { background: var(--eapp-card); border-bottom: 1rpx solid var(--eapp-border); }
.search-bar { background: var(--eapp-card); padding: 12rpx 24rpx 16rpx; border-bottom: 1rpx solid var(--eapp-border); }

.list { padding: 20rpx; display: grid; gap: 16rpx; padding-bottom: 240rpx; }
.fab { position: fixed; right: 24rpx; bottom: calc(140rpx + env(safe-area-inset-bottom)); display: grid; gap: 12rpx; z-index: 20; }
.picker { display: inline-flex; height: 60rpx; align-items: center; padding: 0 18rpx; border: 1rpx solid var(--eapp-border); border-radius: var(--eapp-radius); background: var(--eapp-card); }
.mt { margin-top: 12rpx; }
</style>
