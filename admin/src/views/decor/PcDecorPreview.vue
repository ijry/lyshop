<template>
  <div class="bg-gray-50 min-h-full">
    <template v-for="(comp, index) in components" :key="comp.id">
      <!-- Hero -->
      <section v-if="comp.type === 'hero'" class="relative overflow-hidden cursor-pointer"
        @click.stop="emit('select', index)"
        :style="{ background: `linear-gradient(135deg, ${comp.props?.bg_from || '#b91c1c'}, ${comp.props?.bg_to || '#991b1b'})` }">
        <div class="absolute inset-0" style="background: radial-gradient(circle at 30% 50%, rgba(255,255,255,0.08), transparent 50%)" />
        <div class="max-w-7xl mx-auto px-6 py-16 relative">
          <div class="max-w-2xl">
            <div v-if="comp.props?.badge" class="inline-flex items-center gap-2 bg-white/10 rounded-full px-3 py-1 mb-4">
              <span class="w-1.5 h-1.5 bg-green-400 rounded-full" />
              <span class="text-white/80 text-xs">{{ comp.props.badge }}</span>
            </div>
            <h1 class="text-3xl font-bold text-white leading-tight mb-3">
              <template v-for="(line, i) in splitLines(comp.props?.title)" :key="i">
                <br v-if="i > 0" />{{ line }}
              </template>
            </h1>
            <p v-if="comp.props?.subtitle" class="text-sm mb-6 leading-relaxed" style="color: rgba(255,255,255,0.75)">
              {{ comp.props.subtitle }}
            </p>
            <div class="flex gap-2">
              <span v-if="comp.props?.btn_text"
                class="px-6 py-2 rounded-lg font-semibold text-xs bg-white cursor-default"
                :style="{ color: comp.props.bg_from || '#b91c1c' }">
                {{ comp.props.btn_text }}
              </span>
              <span v-if="comp.props?.btn2_text"
                class="px-6 py-2 rounded-lg font-semibold text-xs text-white border border-white/20 bg-white/10 cursor-default">
                {{ comp.props.btn2_text }}
              </span>
            </div>
          </div>
        </div>
      </section>

      <!-- Banner -->
      <section v-else-if="comp.type === 'banner'" class="max-w-7xl mx-auto px-6 mt-6 cursor-pointer" @click.stop="emit('select', index)">
        <div class="rounded-xl overflow-hidden bg-gray-200" :style="{ height: (comp.props?.height || 400) + 'px' }">
          <img v-if="(comp.props?.images || [])[0]?.url" :src="comp.props.images[0].url"
            class="w-full h-full object-cover" />
          <div v-else class="w-full h-full flex-center text-gray-400 text-sm">轮播图预览</div>
        </div>
      </section>

      <!-- Category Nav -->
      <section v-else-if="comp.type === 'category_nav'"
        @click.stop="emit('select', index)"
        class="cursor-pointer"
        :class="comp.props?.style === 'floating' ? 'max-w-7xl mx-auto px-6 -mt-6 relative z-10' : 'max-w-7xl mx-auto px-6 mt-6'">
        <div :class="comp.props?.style === 'floating' ? 'bg-white rounded-xl shadow-md p-5' : 'bg-white rounded-xl p-5 border border-gray-100'">
          <div class="grid gap-3" :style="{ gridTemplateColumns: `repeat(${comp.props?.columns || 8}, 1fr)` }">
            <div v-for="item in (comp.props?.items || [])" :key="item.title" class="flex flex-col items-center gap-1.5">
              <div class="w-10 h-10 rounded-lg bg-blue-50 flex-center">
                <span class="text-xs font-medium text-blue-600">{{ item.title?.slice(0, 2) }}</span>
              </div>
              <span class="text-xs text-gray-600">{{ item.title }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Product Grid -->
      <section v-else-if="comp.type === 'product_grid'" class="max-w-7xl mx-auto px-6 mt-8 cursor-pointer" @click.stop="emit('select', index)">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-bold text-gray-900">{{ comp.props?.title || '推荐商品' }}</h2>
          <span class="text-xs text-blue-600">查看全部 &rarr;</span>
        </div>
        <div class="grid gap-4" :style="{ gridTemplateColumns: `repeat(${comp.props?.columns || 4}, 1fr)` }">
          <div v-for="product in resolveProductGridItems(comp.props)" :key="product.id"
            class="bg-white rounded-xl border border-gray-100 overflow-hidden">
            <img v-if="product.cover" :src="product.cover" class="aspect-square w-full bg-gray-100 object-cover" />
            <div v-else class="aspect-square bg-gray-100" />
            <div class="p-3">
              <p class="text-sm font-medium text-gray-800 line-clamp-2 min-h-[2.5rem]">{{ product.title }}</p>
              <div class="mt-3 flex items-end justify-between gap-2">
                <span class="text-base font-bold text-red-600">¥{{ formatPrice(product.price) }}</span>
                <span class="text-xs text-gray-400">销量 {{ product.sales || 0 }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Notice -->
      <section v-else-if="comp.type === 'notice'" class="max-w-7xl mx-auto px-6 mt-5 cursor-pointer" @click.stop="emit('select', index)">
        <div class="rounded-lg px-4 py-2.5 flex items-center gap-2 text-sm"
          :style="{ background: comp.props?.bgColor || '#fff7ed', color: comp.props?.color || '#f97316' }">
          <span class="i-carbon-volume-up text-base shrink-0" />
          <span v-for="(item, i) in (comp.props?.items || []).slice(0, 1)" :key="i">{{ item.text }}</span>
        </div>
      </section>

      <!-- Image Ad -->
      <section v-else-if="comp.type === 'image_ad'" class="max-w-7xl mx-auto px-6 mt-6 cursor-pointer" @click.stop="emit('select', index)">
        <div class="rounded-xl overflow-hidden bg-gray-200"
          :style="{ height: comp.props?.height ? comp.props.height + 'px' : '180px' }">
          <img v-if="comp.props?.url" :src="comp.props.url" class="w-full h-full object-cover" />
          <div v-else class="w-full h-full flex-center text-gray-400 text-sm">广告图预览</div>
        </div>
      </section>

      <!-- Rich Text -->
      <section v-else-if="comp.type === 'rich_text'" class="max-w-7xl mx-auto px-6 mt-6 cursor-pointer" @click.stop="emit('select', index)">
        <div class="prose prose-sm max-w-none" v-html="comp.props?.content || '<p class=&quot;text-gray-400&quot;>富文本内容</p>'" />
      </section>

      <!-- Marketing Zone -->
      <section v-else-if="comp.type === 'marketing_zone'" class="max-w-7xl mx-auto px-6 mt-6 cursor-pointer" @click.stop="emit('select', index)">
        <div class="rounded-xl p-5 text-white"
          :style="{ background: `linear-gradient(135deg, ${comp.props?.bg_from || '#b91c1c'}, ${comp.props?.bg_to || '#dc2626'})` }">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <h3 class="text-base font-bold">{{ comp.props?.title || '限时秒杀' }}</h3>
              <span class="text-xs opacity-80">{{ comp.props?.subtitle || '限时抢购中' }}</span>
            </div>
            <span class="text-xs opacity-80">更多 &rarr;</span>
          </div>
          <div v-if="(comp.props?.products || []).length" class="grid grid-cols-4 gap-3">
            <div v-for="p in comp.props.products" :key="p.product_id"
              class="bg-white/10 rounded-lg p-2">
              <img :src="p.cover" class="w-full aspect-square rounded object-cover mb-1.5" />
              <p class="text-xs line-clamp-1 mb-0.5">{{ p.title }}</p>
              <span class="text-xs font-bold">¥{{ p.activity_price || p.group_price || p.floor_price }}</span>
            </div>
          </div>
          <div v-else class="text-center py-6 text-white/50 text-xs">拖入商品数据</div>
        </div>
      </section>

      <!-- Features -->
      <section v-else-if="comp.type === 'features'" class="max-w-7xl mx-auto px-6 mt-8 mb-6 cursor-pointer" @click.stop="emit('select', index)">
        <div class="grid gap-3" :style="{ gridTemplateColumns: `repeat(${comp.props?.columns || 4}, 1fr)` }">
          <div v-for="f in (comp.props?.items || [])" :key="f.title"
            class="flex items-center gap-2.5 bg-white rounded-lg p-3 border border-gray-100">
            <div :class="f.icon" class="text-xl shrink-0 text-blue-600" />
            <div>
              <p class="text-xs font-semibold text-gray-800">{{ f.title }}</p>
              <p class="text-xs text-gray-400">{{ f.desc }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Grid -->
      <section v-else-if="comp.type === 'grid'" class="max-w-7xl mx-auto px-6 mt-6 cursor-pointer" @click.stop="emit('select', index)">
        <div class="bg-white rounded-xl p-5 border border-gray-100">
          <div class="grid gap-3" :style="{ gridTemplateColumns: `repeat(${comp.props?.columns || 4}, 1fr)` }">
            <div v-for="item in (comp.props?.items || [])" :key="item.title"
              class="flex flex-col items-center gap-1.5 py-1.5">
              <div class="w-9 h-9 rounded-lg flex-center text-base" :style="{ background: item.bg || '#f5f5f5' }">
                {{ item.icon }}
              </div>
              <span class="text-xs text-gray-600">{{ item.title }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Spacer -->
      <div v-else-if="comp.type === 'spacer'"
        class="cursor-pointer"
        @click.stop="emit('select', index)"
        :style="{ height: (comp.props?.height || 16) + 'px', background: comp.props?.background || 'transparent' }" />
    </template>

    <!-- Empty state -->
    <div v-if="!components.length" class="flex-center py-20 text-gray-300 text-sm">
      从左侧组件库拖入组件开始装修
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getProducts } from '@/api/plugins'
import productsData from '../../../../app/mock/data/products.json'

defineProps<{ components: any[] }>()

const emit = defineEmits<{
  select: [index: number]
}>()

const previewProducts = ref<any[]>([])

function getMockProducts() {
  const list = Array.isArray((productsData as any)?.list) ? (productsData as any).list : []
  return list.map((item: any) => ({
    ...item,
    price: Number(item?.price || 0),
    sales: Number(item?.sales || 0),
    created_at: item?.created_at || item?.createdAt || '',
  }))
}

function sortProductsBySource(source: string, list: any[]) {
  const rows = [...list]
  if (source === 'new') {
    return rows.sort((left, right) => String(right?.created_at || '').localeCompare(String(left?.created_at || '')))
  }
  if (source === 'recommend') {
    return rows.sort((left, right) => Number(right?.favorite_count || 0) - Number(left?.favorite_count || 0))
  }
  return rows.sort((left, right) => Number(right?.sales || 0) - Number(left?.sales || 0))
}

function resolveProductGridItems(props: any) {
  const source = String(props?.source || 'hot')
  const limit = Math.max(1, Number(props?.limit || 8))
  return sortProductsBySource(source, previewProducts.value).slice(0, limit)
}

function formatPrice(price: number) {
  return Number(price || 0).toFixed(2)
}

onMounted(async () => {
  try {
    const data: any = await getProducts({ page: 1, size: 50 })
    const list = Array.isArray(data?.list) ? data.list : []
    previewProducts.value = list.length ? list : getMockProducts()
  } catch {
    previewProducts.value = getMockProducts()
  }
})

function splitLines(text?: string): string[] {
  return (text || '').split(/\\n|\n/)
}
</script>
