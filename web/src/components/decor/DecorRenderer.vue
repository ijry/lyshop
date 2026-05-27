<template>
  <div>
    <template v-for="comp in components" :key="comp.id">
      <!-- Hero -->
      <section v-if="comp.type === 'hero'" class="relative overflow-hidden"
        :style="{ background: `linear-gradient(135deg, ${comp.props?.bg_from || 'var(--color-hero-from, #b91c1c)'}, ${comp.props?.bg_to || 'var(--color-hero-to, #991b1b)'})` }">
        <div class="absolute inset-0 bg-[radial-gradient(circle_at_30%_50%,rgba(255,255,255,0.08),transparent_50%)]" />
        <div class="max-w-7xl mx-auto px-6 py-20 md:py-28 relative">
          <div class="max-w-2xl">
            <div v-if="comp.props?.badge" class="inline-flex items-center gap-2 bg-white/10 backdrop-blur-sm rounded-full px-4 py-1.5 mb-6">
              <span class="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
              <span class="text-white/80 text-xs font-medium">{{ comp.props.badge }}</span>
            </div>
            <h1 class="text-4xl md:text-5xl font-bold text-white leading-tight mb-4">
              <template v-for="(line, i) in splitLines(comp.props?.title)" :key="i">
                <br v-if="i > 0" />{{ line }}
              </template>
            </h1>
            <p v-if="comp.props?.subtitle" class="text-lg mb-8 leading-relaxed" style="color: rgba(255,255,255,0.75)">
              {{ comp.props.subtitle }}
            </p>
            <div class="flex gap-3">
              <router-link v-if="comp.props?.btn_text" :to="comp.props.btn_link || '/products'"
                class="px-8 py-3 rounded-xl font-semibold text-sm hover:opacity-90 transition-colors"
                :style="{ background: 'white', color: comp.props.bg_from || 'var(--color-primary)' }">
                {{ comp.props.btn_text }}
              </router-link>
              <router-link v-if="comp.props?.btn2_text" :to="comp.props.btn2_link || '/products'"
                class="bg-white/10 backdrop-blur-sm text-white px-8 py-3 rounded-xl font-semibold text-sm border border-white/20 hover:bg-white/20 transition-colors">
                {{ comp.props.btn2_text }}
              </router-link>
            </div>
          </div>
        </div>
      </section>

      <!-- Banner Carousel -->
      <section v-else-if="comp.type === 'banner'" class="max-w-7xl mx-auto px-6 mt-8">
        <div class="relative rounded-2xl overflow-hidden shadow-lg" :style="{ height: (comp.props?.height || 400) + 'px' }">
          <div class="flex transition-transform duration-500"
            :style="{ transform: `translateX(-${(bannerIndex[comp.id] || 0) * 100}%)` }">
            <div v-for="(img, i) in (comp.props?.images || [])" :key="i"
              class="w-full shrink-0 cursor-pointer"
              :style="{ height: (comp.props?.height || 400) + 'px' }"
              @click="navigate(img.link)">
              <img :src="img.url" :alt="img.alt || ''" class="w-full h-full object-cover" />
            </div>
          </div>
          <div v-if="(comp.props?.images || []).length > 1" class="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
            <button v-for="(_, i) in comp.props.images" :key="i"
              class="w-2.5 h-2.5 rounded-full transition-colors"
              :class="(bannerIndex[comp.id] || 0) === i ? 'bg-white' : 'bg-white/40'"
              @click.stop="bannerIndex[comp.id] = i" />
          </div>
        </div>
      </section>

      <!-- Category Nav -->
      <section v-else-if="comp.type === 'category_nav'"
        :class="comp.props?.style === 'floating' ? 'max-w-7xl mx-auto px-6 -mt-8 relative z-10' : 'max-w-7xl mx-auto px-6 mt-8'">
        <div :class="comp.props?.style === 'floating' ? 'bg-white rounded-2xl shadow-lg shadow-gray-200/50 p-6' : 'bg-white rounded-2xl p-6 border border-gray-100'">
          <div class="grid gap-4" :style="{ gridTemplateColumns: `repeat(${comp.props?.columns || 8}, 1fr)` }">
            <div v-for="item in (comp.props?.items || [])" :key="item.title"
              @click="navigate(item.link)" class="flex flex-col items-center gap-2 cursor-pointer group">
              <div v-if="item.icon" class="w-12 h-12 rounded-xl overflow-hidden">
                <img :src="item.icon" class="w-full h-full object-cover" />
              </div>
              <div v-else class="w-12 h-12 rounded-xl flex-center transition-colors"
                :style="{ background: 'color-mix(in srgb, var(--color-primary) 10%, white)', color: 'var(--color-primary)' }">
                <span class="text-sm font-medium">{{ item.title?.slice(0, 2) }}</span>
              </div>
              <span class="text-xs text-gray-600 group-hover:text-[var(--color-primary)] transition-colors">{{ item.title }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Product Grid -->
      <section v-else-if="comp.type === 'product_grid'" class="max-w-7xl mx-auto px-6 mt-12">
        <div class="flex-between mb-6">
          <h2 class="text-xl font-bold text-gray-900">{{ comp.props?.title || '推荐商品' }}</h2>
          <router-link to="/products"
            class="text-sm font-medium flex items-center gap-1 hover:opacity-80 transition-colors"
            :style="{ color: 'var(--color-primary)' }">
            查看全部 <div class="i-carbon-arrow-right text-sm" />
          </router-link>
        </div>
        <div class="grid gap-5" :style="{ gridTemplateColumns: `repeat(${comp.props?.columns || 4}, 1fr)` }">
          <ProductCard v-for="p in (gridProducts[comp.id] || [])" :key="p.id" :product="p" />
        </div>
      </section>

      <!-- Notice -->
      <section v-else-if="comp.type === 'notice'" class="max-w-7xl mx-auto px-6 mt-6">
        <div class="rounded-xl px-5 py-3 flex items-center gap-3 overflow-hidden"
          :style="{ background: comp.props?.bgColor || '#fff7ed', color: comp.props?.color || '#f97316' }">
          <div class="i-carbon-volume-up text-lg shrink-0" />
          <div class="overflow-hidden flex-1 relative h-6">
            <div class="notice-scroll absolute whitespace-nowrap"
              :style="{ animationDuration: ((comp.props?.items || []).length * 4) + 's' }">
              <span v-for="(item, i) in (comp.props?.items || [])" :key="i" class="mr-12 text-sm cursor-pointer hover:underline"
                @click="navigate(item.link)">{{ item.text }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Image Ad -->
      <section v-else-if="comp.type === 'image_ad'" class="max-w-7xl mx-auto px-6 mt-8">
        <div class="rounded-2xl overflow-hidden cursor-pointer hover:shadow-lg transition-shadow"
          @click="navigate(comp.props?.link)">
          <img v-if="comp.props?.url" :src="comp.props.url" :alt="comp.props?.alt || ''"
            class="w-full block" :style="{ height: comp.props?.height ? comp.props.height + 'px' : 'auto' }" />
        </div>
      </section>

      <!-- Rich Text -->
      <section v-else-if="comp.type === 'rich_text'" class="max-w-7xl mx-auto px-6 mt-8">
        <div class="prose max-w-none" v-html="comp.props?.content || ''" />
      </section>

      <!-- Marketing Zone -->
      <section v-else-if="comp.type === 'marketing_zone'" class="max-w-7xl mx-auto px-6 mt-8">
        <div class="rounded-2xl p-6 text-white"
          :style="{ background: `linear-gradient(135deg, ${comp.props?.bg_from || 'var(--color-primary-dark)'}, ${comp.props?.bg_to || 'var(--color-primary)'})` }">
          <div class="flex-between mb-4">
            <div class="flex items-center gap-3">
              <h3 class="text-lg font-bold">{{ comp.props?.title || '限时秒杀' }}</h3>
              <span class="text-sm opacity-80">{{ comp.props?.subtitle || '限时抢购中' }}</span>
            </div>
            <router-link :to="comp.props?.more_link || '/products'" class="text-sm opacity-80 hover:opacity-100 transition-opacity">
              更多 <span class="i-carbon-arrow-right inline-block align-middle" />
            </router-link>
          </div>
          <div v-if="(comp.props?.products || []).length" class="grid grid-cols-4 gap-4">
            <div v-for="p in comp.props.products" :key="p.product_id"
              @click="$router.push(`/product/${p.product_id}`)"
              class="bg-white/10 backdrop-blur-sm rounded-xl p-3 cursor-pointer hover:bg-white/20 transition-colors">
              <img :src="p.cover" class="w-full aspect-square rounded-lg object-cover mb-2" />
              <p class="text-xs line-clamp-1 mb-1">{{ p.title }}</p>
              <div class="flex items-baseline gap-2">
                <span class="text-sm font-bold">¥{{ p.activity_price || p.group_price || p.floor_price }}</span>
                <span class="text-xs line-through opacity-60">¥{{ p.origin_price }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Features -->
      <section v-else-if="comp.type === 'features'" class="max-w-7xl mx-auto px-6 mt-12 mb-8">
        <div class="grid gap-4" :style="{ gridTemplateColumns: `repeat(${comp.props?.columns || 4}, 1fr)` }">
          <div v-for="f in (comp.props?.items || [])" :key="f.title"
            class="flex items-center gap-3 bg-white rounded-xl p-4 border border-gray-100">
            <div :class="f.icon" class="text-2xl shrink-0" :style="{ color: 'var(--color-primary)' }" />
            <div>
              <p class="text-sm font-semibold text-gray-800">{{ f.title }}</p>
              <p class="text-xs text-gray-400 mt-0.5">{{ f.desc }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Grid (quick entry icons) -->
      <section v-else-if="comp.type === 'grid'" class="max-w-7xl mx-auto px-6 mt-8">
        <div class="bg-white rounded-2xl p-6 border border-gray-100">
          <div class="grid gap-4" :style="{ gridTemplateColumns: `repeat(${comp.props?.columns || 4}, 1fr)` }">
            <div v-for="item in (comp.props?.items || [])" :key="item.title"
              @click="navigate(item.link)" class="flex flex-col items-center gap-2 cursor-pointer group py-2">
              <div class="w-11 h-11 rounded-xl flex-center text-xl"
                :style="{ background: item.bg || '#f5f5f5' }">
                {{ item.icon }}
              </div>
              <span class="text-xs text-gray-600 group-hover:text-[var(--color-primary)] transition-colors">{{ item.title }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Spacer -->
      <div v-else-if="comp.type === 'spacer'"
        :style="{ height: (comp.props?.height || 16) + 'px', background: comp.props?.background || 'transparent' }" />
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { get } from '@/api/request'
import ProductCard from '@/components/ProductCard.vue'

const props = defineProps<{ components: any[] }>()
const router = useRouter()
const gridProducts = ref<Record<string, any[]>>({})
const bannerIndex = ref<Record<string, number>>({})
let bannerTimers: Record<string, ReturnType<typeof setInterval>> = {}

function navigate(link?: string) {
  if (!link) return
  if (link.startsWith('http')) {
    window.open(link, '_blank')
  } else {
    const pcLink = link
      .replace(/^\/pages\/product\/list/, '/products')
      .replace(/^\/pages\/product\/detail\?id=/, '/product/')
      .replace(/^\/pages\/marketing\/seckill$/, '/products/seckill')
      .replace(/^\/pages\/marketing\/group-buy$/, '/products/group-buy')
      .replace(/^\/pages\/marketing\/bargain$/, '/products/bargain')
      .replace(/^\/pages\//, '/')
    router.push(pcLink)
  }
}

function splitLines(text?: string): string[] {
  return (text || '').split(/\\n|\n/)
}

async function loadGridProducts() {
  const next: Record<string, any[]> = {}
  await Promise.all(
    props.components
      .filter(c => c.type === 'product_grid')
      .map(async (comp) => {
        const data = await get<any>('/api/v1/products', { page: 1, size: comp.props?.limit || 8 })
        next[comp.id] = data?.list || []
      })
  )
  gridProducts.value = next
}

function startBannerAutoplay() {
  stopBannerAutoplay()
  props.components
    .filter(c => c.type === 'banner' && (c.props?.images || []).length > 1)
    .forEach(comp => {
      bannerIndex.value[comp.id] = 0
      bannerTimers[comp.id] = setInterval(() => {
        const len = comp.props.images.length
        bannerIndex.value[comp.id] = ((bannerIndex.value[comp.id] || 0) + 1) % len
      }, comp.props?.interval || 4000)
    })
}

function stopBannerAutoplay() {
  Object.values(bannerTimers).forEach(clearInterval)
  bannerTimers = {}
}

onMounted(() => {
  loadGridProducts()
  startBannerAutoplay()
})
onUnmounted(stopBannerAutoplay)
watch(() => props.components, () => {
  loadGridProducts()
  startBannerAutoplay()
}, { deep: true })
</script>

<style scoped>
@keyframes notice-scroll {
  0% { transform: translateX(100%); }
  100% { transform: translateX(-100%); }
}
.notice-scroll {
  animation: notice-scroll 12s linear infinite;
}
</style>
