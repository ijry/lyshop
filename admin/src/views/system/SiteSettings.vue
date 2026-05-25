<template>
  <div>
    <h2 class="text-xl font-semibold text-slate-800 mb-6">站点设置</h2>

    <div v-if="loading" class="text-center py-12 text-slate-400">加载中...</div>

    <div v-else class="flex gap-6">
      <!-- Section tabs -->
      <div class="w-48 shrink-0">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden">
          <button v-for="sec in sections" :key="sec.key"
            @click="activeSection = sec.key"
            :class="activeSection === sec.key ? 'bg-red-50 text-red-600 font-medium border-l-3 border-l-red-600' : 'text-slate-600 hover:bg-slate-50'"
            class="w-full text-left px-4 py-3 text-sm transition-colors border-b border-slate-50 last:border-0">
            {{ sec.title }}
          </button>
        </div>
      </div>

      <!-- Form -->
      <div class="flex-1 max-w-2xl">
        <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6">

          <!-- 基本信息 -->
          <div v-show="activeSection === 'basic'" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">站点名称</label>
              <input v-model="form.site_name" type="text" placeholder="LYShop"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">站点Logo URL</label>
              <input v-model="form.site_logo" type="text" placeholder="留空则使用站点名首字母"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">SEO 标题</label>
              <input v-model="form.seo_title" type="text" placeholder="页面 title 标签"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">SEO 关键词</label>
              <input v-model="form.seo_keywords" type="text" placeholder="多个关键词用英文逗号分隔"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">SEO 描述</label>
              <textarea v-model="form.seo_description" rows="3" placeholder="搜索引擎描述"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm resize-y focus:outline-none focus:border-red-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">ICP 备案号</label>
              <input v-model="form.icp" type="text" placeholder="如：京ICP备XXXXXXXX号"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
            </div>
          </div>

          <!-- Hero 配置 -->
          <div v-show="activeSection === 'hero'" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">角标文字</label>
              <input v-model="form.hero_badge" type="text" placeholder="限时秒杀进行中"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">主标题</label>
              <textarea v-model="form.hero_title" rows="2" placeholder="精选好物\n品质生活从这里开始"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm resize-y focus:outline-none focus:border-red-500" />
              <p class="text-xs text-slate-400 mt-1">用 \n 换行</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">副标题</label>
              <input v-model="form.hero_subtitle" type="text"
                class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1.5">按钮文字</label>
                <input v-model="form.hero_btn_text" type="text" placeholder="立即选购"
                  class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1.5">按钮链接</label>
                <input v-model="form.hero_btn_link" type="text" placeholder="/products"
                  class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
              </div>
            </div>
            <!-- Preview -->
            <div class="mt-4 rounded-xl overflow-hidden">
              <div class="p-6 text-white text-center"
                :style="{ background: `linear-gradient(135deg, ${form.color_hero_from}, ${form.color_hero_to})` }">
                <div class="text-xs opacity-80 mb-2">{{ form.hero_badge }}</div>
                <div class="text-lg font-bold mb-1" v-html="form.hero_title.replace(/\\n/g, '<br>')"></div>
                <div class="text-sm opacity-80">{{ form.hero_subtitle }}</div>
              </div>
            </div>
          </div>

          <!-- 主题色 -->
          <div v-show="activeSection === 'theme'" class="space-y-5">
            <p class="text-sm text-slate-500 mb-2">配置 PC 端主题色，保存后前端实时生效。</p>
            <div class="grid grid-cols-2 gap-x-8 gap-y-4">
              <ColorInput v-model="form.color_primary" label="主色" />
              <ColorInput v-model="form.color_primary_light" label="主色（浅）" />
              <ColorInput v-model="form.color_primary_dark" label="主色（深）" />
              <ColorInput v-model="form.color_price" label="价格色" />
              <ColorInput v-model="form.color_hero_from" label="Hero渐变起" />
              <ColorInput v-model="form.color_hero_to" label="Hero渐变止" />
              <ColorInput v-model="form.color_bg_page" label="页面背景" />
              <ColorInput v-model="form.color_bg_footer" label="底栏背景" />
            </div>
            <!-- Theme preview -->
            <div class="mt-4 p-4 rounded-xl border border-slate-200" :style="{ background: form.color_bg_page }">
              <div class="flex items-center gap-3 mb-3">
                <div class="w-8 h-8 rounded-lg flex items-center justify-center text-white text-sm font-bold" :style="{ background: form.color_primary }">L</div>
                <span class="font-bold text-slate-800">{{ form.site_name || 'LYShop' }}</span>
              </div>
              <div class="h-16 rounded-lg mb-3" :style="{ background: `linear-gradient(135deg, ${form.color_hero_from}, ${form.color_hero_to})` }"></div>
              <div class="flex gap-4">
                <span class="text-sm font-bold" :style="{ color: form.color_price }">¥99.00</span>
                <span class="text-sm" :style="{ color: form.color_primary }">立即购买</span>
              </div>
            </div>
          </div>

          <!-- Save -->
          <div class="flex items-center gap-3 mt-6 pt-4 border-t border-slate-100">
            <button @click="save" :disabled="saving"
              class="px-6 py-2.5 bg-red-600 text-white rounded-xl text-sm font-medium hover:bg-red-700 transition disabled:opacity-40">
              {{ saving ? '保存中...' : '保存设置' }}
            </button>
            <span v-if="savedMsg" class="text-sm text-green-600">{{ savedMsg }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import request from '@/api/request'
import ColorInput from '@/views/decor/widgets/ColorInput.vue'

const sections = [
  { key: 'basic', title: '基本信息' },
  { key: 'hero', title: '首页 Hero' },
  { key: 'theme', title: '主题色（PC）' },
]

const activeSection = ref('basic')
const loading = ref(true)
const saving = ref(false)
const savedMsg = ref('')

const form = reactive({
  site_name: 'LYShop',
  site_logo: '',
  seo_title: 'LYShop - 开源商城',
  seo_keywords: '商城,电商,开源',
  seo_description: '开源插件化商城系统',
  icp: '',
  hero_badge: '限时秒杀进行中',
  hero_title: '精选好物\\n品质生活从这里开始',
  hero_subtitle: '数千款精选商品，正品保障，极速发货，让购物更简单。',
  hero_btn_text: '立即选购',
  hero_btn_link: '/products',
  color_primary: '#dc2626',
  color_primary_light: '#ef4444',
  color_primary_dark: '#b91c1c',
  color_bg_page: '#f9fafb',
  color_bg_header: 'rgba(255,255,255,0.8)',
  color_bg_footer: '#f9fafb',
  color_price: '#ef4444',
  color_hero_from: '#b91c1c',
  color_hero_to: '#991b1b',
})

async function save() {
  saving.value = true
  try {
    await request.put('/site-settings', { ...form })
    savedMsg.value = '保存成功'
    setTimeout(() => savedMsg.value = '', 2000)
  } catch (e: any) {
    savedMsg.value = '保存失败: ' + (e.message || '')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const data: any = await request.get('/site-settings')
    if (data) Object.assign(form, data)
  } finally {
    loading.value = false
  }
})
</script>
