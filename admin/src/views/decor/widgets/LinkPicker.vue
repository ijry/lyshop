<template>
  <div class="relative">
    <label v-if="label" class="block text-xs text-slate-500 mb-1">{{ label }}</label>
    <div class="flex items-center gap-1.5">
      <div class="flex-1 border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs text-slate-600 truncate bg-slate-50 min-w-0">
        {{ displayText }}
      </div>
      <button @click="open = !open"
        class="px-2.5 py-1.5 text-xs bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition shrink-0">
        选择
      </button>
    </div>

    <!-- Popover -->
    <div v-if="open"
      class="absolute left-0 right-0 top-full mt-1 z-50 bg-white rounded-xl shadow-xl border border-slate-200 overflow-hidden"
      style="max-height: 380px;">
      <!-- Tabs -->
      <div class="flex border-b border-slate-100 text-xs">
        <button v-for="tab in tabs" :key="tab.key" @click="activeTab = tab.key"
          :class="activeTab === tab.key ? 'text-blue-600 border-b-2 border-blue-600' : 'text-slate-500'"
          class="flex-1 px-2 py-2 text-center hover:bg-slate-50 transition">
          {{ tab.label }}
        </button>
      </div>

      <div class="overflow-y-auto" style="max-height: 330px;">
        <!-- Tab: 商城页面 -->
        <div v-if="activeTab === 'pages'" class="p-2 space-y-0.5">
          <div v-for="page in shopPages" :key="page.path" @click="selectLink(page.path)"
            :class="modelValue === page.path ? 'bg-blue-50 text-blue-600' : 'hover:bg-slate-50'"
            class="px-3 py-2 rounded-lg text-xs cursor-pointer transition">
            {{ page.label }}
            <span class="text-slate-400 ml-1">{{ page.path }}</span>
          </div>
        </div>

        <!-- Tab: 商品分类 -->
        <div v-if="activeTab === 'categories'" class="p-2">
          <div v-if="!categories.length" class="text-center py-4 text-xs text-slate-400">加载中...</div>
          <div v-for="cat in categories" :key="cat.id" @click="selectLink(`/pages/product/list?category_id=${cat.id}`)"
            :class="modelValue === `/pages/product/list?category_id=${cat.id}` ? 'bg-blue-50 text-blue-600' : 'hover:bg-slate-50'"
            class="px-3 py-2 rounded-lg text-xs cursor-pointer transition">
            {{ cat.name }}
          </div>
        </div>

        <!-- Tab: 指定商品 -->
        <div v-if="activeTab === 'products'" class="p-2 space-y-2">
          <input v-model="productKeyword" @input="searchProducts" placeholder="搜索商品名称..."
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
          <div v-if="!products.length" class="text-center py-3 text-xs text-slate-400">
            {{ productKeyword ? '未找到商品' : '输入关键词搜索' }}
          </div>
          <div v-for="p in products" :key="p.id" @click="selectLink(`/pages/product/detail?id=${p.id}`)"
            :class="modelValue === `/pages/product/detail?id=${p.id}` ? 'bg-blue-50' : 'hover:bg-slate-50'"
            class="flex items-center gap-2 px-2 py-1.5 rounded-lg cursor-pointer transition">
            <img v-if="p.cover" :src="p.cover" class="w-8 h-8 rounded object-cover shrink-0" />
            <div v-else class="w-8 h-8 rounded bg-slate-100 shrink-0" />
            <div class="min-w-0">
              <div class="text-xs text-slate-700 truncate">{{ p.title }}</div>
              <div class="text-xs text-red-500">¥{{ p.price }}</div>
            </div>
          </div>
        </div>

        <!-- Tab: 营销活动 -->
        <div v-if="activeTab === 'marketing'" class="p-2 space-y-0.5">
          <div v-for="item in marketingPages" :key="item.path" @click="selectLink(item.path)"
            :class="modelValue === item.path ? 'bg-blue-50 text-blue-600' : 'hover:bg-slate-50'"
            class="px-3 py-2 rounded-lg text-xs cursor-pointer transition">
            {{ item.label }}
          </div>
        </div>

        <!-- Tab: 自定义 -->
        <div v-if="activeTab === 'custom'" class="p-3 space-y-2">
          <input v-model="customUrl" placeholder="/pages/..."
            class="w-full border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs focus:outline-none focus:border-blue-400" />
          <button @click="selectLink(customUrl)" :disabled="!customUrl.trim()"
            class="w-full px-3 py-1.5 text-xs bg-blue-600 text-white rounded-lg hover:bg-blue-500 disabled:opacity-50 transition">
            确认
          </button>
        </div>
      </div>
    </div>

    <!-- Backdrop -->
    <div v-if="open" class="fixed inset-0 z-40" @click="open = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { getCategories, getProducts } from '@/api/plugins'

const props = defineProps<{
  modelValue: string
  label?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const open = ref(false)
const activeTab = ref<'pages' | 'categories' | 'products' | 'marketing' | 'custom'>('pages')
const categories = ref<any[]>([])
const products = ref<any[]>([])
const productKeyword = ref('')
const customUrl = ref('')

let categoriesLoaded = false
let searchTimer: ReturnType<typeof setTimeout> | null = null

const tabs = [
  { key: 'pages' as const, label: '页面' },
  { key: 'categories' as const, label: '分类' },
  { key: 'products' as const, label: '商品' },
  { key: 'marketing' as const, label: '营销' },
  { key: 'custom' as const, label: '自定义' },
]

const shopPages = [
  { label: '首页', path: '/pages/index/index' },
  { label: '商品列表', path: '/pages/product/list' },
  { label: '购物车', path: '/pages/cart/index' },
  { label: '我的订单', path: '/pages/order/list' },
  { label: '个人中心', path: '/pages/user/index' },
  { label: '我的收藏', path: '/pages/user/favorites' },
  { label: '我的积分', path: '/pages/user/points' },
  { label: '会员中心', path: '/pages/user/vip' },
  { label: '消息中心', path: '/pages/message/index' },
  { label: '每日签到', path: '/pages/checkin/index' },
]

const marketingPages = [
  { label: '领券中心', path: '/pages/marketing/coupon?mode=claim' },
  { label: '限时秒杀', path: '/pages/marketing/seckill' },
  { label: '拼团特惠', path: '/pages/marketing/group-buy' },
  { label: '砍价免费拿', path: '/pages/marketing/bargain' },
]

const allKnownPages = [...shopPages, ...marketingPages]

const displayText = computed(() => {
  if (!props.modelValue) return '未设置链接'
  const known = allKnownPages.find(p => p.path === props.modelValue)
  if (known) return known.label
  if (props.modelValue.includes('category_id=')) {
    const catId = new URLSearchParams(props.modelValue.split('?')[1]).get('category_id')
    const cat = categories.value.find(c => String(c.id) === catId)
    return cat ? `分类: ${cat.name}` : props.modelValue
  }
  if (props.modelValue.includes('/product/detail')) {
    return `商品详情 ${props.modelValue}`
  }
  return props.modelValue
})

function selectLink(link: string) {
  emit('update:modelValue', link)
  open.value = false
}

watch(activeTab, async (tab) => {
  if (tab === 'categories' && !categoriesLoaded) {
    categoriesLoaded = true
    categories.value = ((await getCategories()) || []) as any[]
  }
  if (tab === 'custom') {
    customUrl.value = props.modelValue || ''
  }
})

function searchProducts() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(async () => {
    if (!productKeyword.value.trim()) {
      products.value = []
      return
    }
    const data: any = await getProducts({ keyword: productKeyword.value, page: 1, size: 20 })
    products.value = data?.list || []
  }, 300)
}
</script>
