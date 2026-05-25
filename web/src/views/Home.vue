<template>
  <div>
    <!-- Hero banner -->
    <section class="relative overflow-hidden"
      :style="{ background: `linear-gradient(135deg, var(--color-hero-from, #b91c1c), var(--color-hero-to, #991b1b))` }">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_30%_50%,rgba(255,255,255,0.08),transparent_50%)]" />
      <div class="max-w-7xl mx-auto px-6 py-20 md:py-28 relative">
        <div class="max-w-2xl">
          <div class="inline-flex items-center gap-2 bg-white/10 backdrop-blur-sm rounded-full px-4 py-1.5 mb-6">
            <span class="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
            <span class="text-white/80 text-xs font-medium">{{ site.settings.hero_badge }}</span>
          </div>
          <h1 class="text-4xl md:text-5xl font-bold text-white leading-tight mb-4">
            <template v-for="(line, i) in heroLines" :key="i">
              <br v-if="i > 0" />{{ line }}
            </template>
          </h1>
          <p class="text-lg mb-8 leading-relaxed" style="color: rgba(255,255,255,0.75)">
            {{ site.settings.hero_subtitle }}
          </p>
          <div class="flex gap-3">
            <router-link :to="site.settings.hero_btn_link || '/products'"
              class="px-8 py-3 rounded-xl font-semibold text-sm hover:opacity-90 transition-colors"
              :style="{ background: 'white', color: 'var(--color-primary, #dc2626)' }">
              {{ site.settings.hero_btn_text }}
            </router-link>
            <router-link to="/products"
              class="bg-white/10 backdrop-blur-sm text-white px-8 py-3 rounded-xl font-semibold text-sm border border-white/20 hover:bg-white/20 transition-colors">
              查看全部
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Categories -->
    <section class="max-w-7xl mx-auto px-6 -mt-8 relative z-10">
      <div class="bg-white rounded-2xl shadow-lg shadow-gray-200/50 p-6">
        <div class="grid grid-cols-4 md:grid-cols-8 gap-4">
          <div v-for="cat in categories" :key="cat.id"
            @click="$router.push(`/products?category=${cat.id}`)"
            class="flex flex-col items-center gap-2 cursor-pointer group">
            <div class="w-12 h-12 rounded-xl flex-center transition-colors"
              :style="{ background: 'color-mix(in srgb, var(--color-primary) 10%, white)', color: 'var(--color-primary)' }">
              <span class="text-sm font-medium">{{ cat.name.slice(0, 2) }}</span>
            </div>
            <span class="text-xs text-gray-600 group-hover:text-[var(--color-primary)] transition-colors">{{ cat.name }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Hot products -->
    <section class="max-w-7xl mx-auto px-6 mt-12">
      <div class="flex-between mb-6">
        <h2 class="text-xl font-bold text-gray-900">热销推荐</h2>
        <router-link to="/products"
          class="text-sm font-medium flex items-center gap-1 hover:opacity-80 transition-colors"
          :style="{ color: 'var(--color-primary, #dc2626)' }">
          查看全部 <div class="i-carbon-arrow-right text-sm" />
        </router-link>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-5">
        <ProductCard v-for="p in products" :key="p.id" :product="p" />
      </div>
    </section>

    <!-- Features -->
    <section class="max-w-7xl mx-auto px-6 mt-16 mb-8">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div v-for="f in features" :key="f.title"
          class="flex items-center gap-3 bg-white rounded-xl p-4 border border-gray-100">
          <div :class="f.icon" class="text-2xl shrink-0" :style="{ color: 'var(--color-primary, #dc2626)' }" />
          <div>
            <p class="text-sm font-semibold text-gray-800">{{ f.title }}</p>
            <p class="text-xs text-gray-400 mt-0.5">{{ f.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <PresetSwitcher />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { get } from '@/api/request'
import { useSiteStore } from '@/stores/site'
import ProductCard from '@/components/ProductCard.vue'
import PresetSwitcher from '@/components/PresetSwitcher.vue'

const site = useSiteStore()
const categories = ref<any[]>([])
const products = ref<any[]>([])

const heroLines = computed(() =>
  (site.settings.hero_title || '').split(/\\n|\n/)
)

const features = [
  { icon: 'i-carbon-delivery-truck', title: '极速发货', desc: '下单24小时内发出' },
  { icon: 'i-carbon-checkmark-outline', title: '正品保障', desc: '全球品牌直供' },
  { icon: 'i-carbon-renew', title: '无忧退换', desc: '7天无理由退货' },
  { icon: 'i-carbon-headset', title: '在线客服', desc: '7×24小时在线' },
]

onMounted(async () => {
  const cats = await get<any[]>('/api/v1/categories')
  categories.value = (cats || []).slice(0, 8)
  const data = await get<any>('/api/v1/products')
  products.value = (data?.list || []).slice(0, 8)
})
</script>
