<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import SkuMatrixEditor from '@/components/biz/SkuMatrixEditor.vue'
import CategoryTreePicker from '@/components/biz/CategoryTreePicker.vue'
import RichTextEditor from '@/components/biz/RichTextEditor.vue'
import { createProduct, getProductDetail, updateProduct } from '@/api/product'

const id = ref(0)
const saving = ref(false)
const showCatPicker = ref(false)

const form = reactive<any>({
  title: '', subtitle: '', sell_points: [] as string[],
  covers: [] as string[],
  detail_html: '',
  price: 0, stock: 0, unit: '件',
  weight: 0,
  category_id: 0, category_path_name: '',
  tags: [] as string[],
  skus: [] as Array<{ id?: number; attrs: Array<{ name: string; value: string }>; price: number; stock: number }>,
  low_stock_threshold: 10,
  shipping_template: 'default',
  limit_per_order: 0,
  exclude_marketing: false,
  status: 1,
  online_at: '', offline_at: '',
})

const newSellPoint = ref('')
const newCover = ref('')
const newTag = ref('')

async function loadData() {
  if (!id.value) return
  const data: any = await getProductDetail(id.value)
  Object.assign(form, {
    title: data?.title || '', subtitle: data?.subtitle || '',
    sell_points: Array.isArray(data?.sell_points) ? data.sell_points : [],
    covers: Array.isArray(data?.covers) ? data.covers : (data?.cover ? [data.cover] : []),
    detail_html: String(data?.detail_html || ''),
    price: Number(data?.price || 0), stock: Number(data?.stock || 0),
    unit: String(data?.unit || '件'), weight: Number(data?.weight || 0),
    category_id: Number(data?.category_id || 0),
    category_path_name: String(data?.category_path_name || ''),
    tags: Array.isArray(data?.tags) ? data.tags : [],
    skus: Array.isArray(data?.skus) ? data.skus : [],
    low_stock_threshold: Number(data?.low_stock_threshold || 10),
    shipping_template: String(data?.shipping_template || 'default'),
    limit_per_order: Number(data?.limit_per_order || 0),
    exclude_marketing: !!data?.exclude_marketing,
    status: Number(data?.status || 0) === 1 ? 1 : 0,
    online_at: String(data?.online_at || ''),
    offline_at: String(data?.offline_at || ''),
  })
}

function addSellPoint() { if (newSellPoint.value.trim()) { form.sell_points.push(newSellPoint.value.trim()); newSellPoint.value = '' } }
function removeSellPoint(i: number) { form.sell_points.splice(i, 1) }
function addCover() { if (newCover.value.trim()) { form.covers.push(newCover.value.trim()); newCover.value = '' } }
function removeCover(i: number) { form.covers.splice(i, 1) }
function addTag() { if (newTag.value.trim()) { form.tags.push(newTag.value.trim()); newTag.value = '' } }
function removeTag(i: number) { form.tags.splice(i, 1) }
function onSkusChange(rows: any[]) { form.skus = rows }
function onPickCategory(payload: { id: number; path_name: string }) {
  form.category_id = payload.id; form.category_path_name = payload.path_name; showCatPicker.value = false
}
function requestEditDetail() {
  uni.setClipboardData({ data: `${form.title} - 详情编辑` })
  uni.showModal({ title: '提示', content: '长详情编辑请到管理后台进行；标题已复制到剪贴板。', showCancel: false })
}

async function save() {
  if (!form.title.trim()) { uni.showToast({ title: '请输入商品标题', icon: 'none' }); return }
  saving.value = true
  try {
    const payload: any = {
      product: {
        title: form.title.trim(), subtitle: form.subtitle.trim(),
        sell_points: form.sell_points,
        cover: form.covers[0] || '',
        covers: form.covers,
        detail_html: form.detail_html,
        price: Number(form.price), stock: Number(form.stock),
        unit: form.unit, weight: Number(form.weight),
        category_id: Number(form.category_id), category_path_name: form.category_path_name,
        tags: form.tags,
        low_stock_threshold: Number(form.low_stock_threshold),
        shipping_template: form.shipping_template,
        limit_per_order: Number(form.limit_per_order),
        exclude_marketing: !!form.exclude_marketing,
        status: form.status,
        online_at: form.online_at, offline_at: form.offline_at,
      },
      skus: form.skus.map((s: any) => ({ ...s, attrs: typeof s.attrs === 'string' ? s.attrs : JSON.stringify(s.attrs), price: Number(s.price), stock: Number(s.stock) })),
    }
    if (id.value) await updateProduct(id.value, payload)
    else await createProduct(payload)
    uni.showToast({ title: '保存成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 350)
  } finally { saving.value = false }
}

onLoad((opts) => { id.value = Number(opts?.id || 0); loadData() })
</script>

<template>
  <view class="page">
    <view class="section">
      <view class="section-title">基础信息</view>
      <up-input v-model="form.title" placeholder="商品标题" class="mt" />
      <up-input v-model="form.subtitle" placeholder="副标题" class="mt" />
      <view class="tag-row mt">
        <view v-for="(sp, i) in form.sell_points" :key="i" class="tag">{{ sp }}<text class="x" @click="removeSellPoint(i)">✕</text></view>
        <up-input v-model="newSellPoint" placeholder="添加卖点" class="tag-input" @confirm="addSellPoint" />
      </view>
    </view>

    <view class="section">
      <view class="section-title">主图轮播</view>
      <view class="cover-row">
        <view v-for="(c, i) in form.covers" :key="i" class="cover-cell">
          <image :src="c" mode="aspectFill" class="cover-img" />
          <text class="x-abs" @click="removeCover(i)">✕</text>
        </view>
      </view>
      <view class="add-cover mt">
        <up-input v-model="newCover" placeholder="图片 URL" />
        <up-button size="mini" type="primary" @click="addCover">添加</up-button>
      </view>
    </view>

    <view class="section">
      <view class="section-title">商品详情</view>
      <RichTextEditor :html="form.detail_html" @requestEdit="requestEditDetail" />
    </view>

    <view class="section">
      <view class="section-title">价格库存</view>
      <view class="grid-2 mt">
        <up-input v-model="form.price" type="digit" placeholder="主价格" />
        <up-input v-model="form.stock" type="number" placeholder="主库存" />
        <up-input v-model="form.unit" placeholder="单位" />
        <up-input v-model="form.weight" type="digit" placeholder="重量(kg)" />
        <up-input v-model="form.low_stock_threshold" type="number" placeholder="预警阈值" />
        <up-input v-model="form.limit_per_order" type="number" placeholder="单笔限购(0=不限)" />
      </view>
    </view>

    <view class="section">
      <view class="section-title">规格 SKU</view>
      <SkuMatrixEditor :skus="form.skus" :base-price="form.price" @update="onSkusChange" />
    </view>

    <view class="section">
      <view class="section-title">分类与标签</view>
      <view class="cat" @click="showCatPicker = true">分类：{{ form.category_path_name || '请选择' }} <text class="caret">›</text></view>
      <view class="tag-row mt">
        <view v-for="(tag, i) in form.tags" :key="i" class="tag">{{ tag }}<text class="x" @click="removeTag(i)">✕</text></view>
        <up-input v-model="newTag" placeholder="添加标签" class="tag-input" @confirm="addTag" />
      </view>
    </view>

    <view class="section">
      <view class="section-title">物流与营销</view>
      <picker mode="selector" :range="['默认模板','包邮','到付','同城']" @change="(e) => form.shipping_template = ['default','free','cod','local'][Number(e.detail.value)]">
        <view class="picker">物流模板：{{ { default: '默认模板', free: '包邮', cod: '到付', local: '同城' }[form.shipping_template] }}</view>
      </picker>
      <view class="row mt">
        <text>不参与营销活动</text>
        <switch :checked="form.exclude_marketing" @change="(e) => form.exclude_marketing = e.detail.value" />
      </view>
    </view>

    <view class="section">
      <view class="section-title">状态控制</view>
      <view class="row">
        <text>上架</text>
        <switch :checked="form.status === 1" @change="(e) => form.status = e.detail.value ? 1 : 0" />
      </view>
      <up-input v-model="form.online_at" placeholder="上架时间 ISO（可选）" class="mt" />
      <up-input v-model="form.offline_at" placeholder="下架时间 ISO（可选）" class="mt" />
    </view>

    <up-button type="primary" :loading="saving" class="save" @click="save">保存</up-button>

    <CategoryTreePicker :show="showCatPicker" :value="form.category_id" @close="showCatPicker = false" @pick="onPickCategory" />
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; padding-bottom: 160rpx; display: grid; gap: 16rpx; box-sizing: border-box; }
.section { background: var(--eapp-card); border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 20rpx; }
.section-title { font-size: 30rpx; font-weight: 700; margin-bottom: 10rpx; display: flex; align-items: center; justify-content: space-between; }
.mt { margin-top: 12rpx; }
.tag-row { display: flex; gap: 10rpx; flex-wrap: wrap; align-items: center; }
.tag { background: var(--eapp-bg); border-radius: 999rpx; padding: 6rpx 16rpx; font-size: 22rpx; display: inline-flex; align-items: center; gap: 6rpx; }
.tag .x { color: var(--eapp-danger); font-size: 20rpx; }
.tag-input { flex: 1; min-width: 180rpx; }
.cover-row { display: flex; gap: 10rpx; flex-wrap: wrap; }
.cover-cell { position: relative; width: 160rpx; height: 160rpx; }
.cover-img { width: 100%; height: 100%; border-radius: 14rpx; }
.x-abs { position: absolute; top: -10rpx; right: -10rpx; background: var(--eapp-danger); color: #fff; font-size: 20rpx; width: 32rpx; height: 32rpx; display: flex; align-items: center; justify-content: center; border-radius: 50%; }
.add-cover { display: flex; gap: 10rpx; align-items: center; }
.grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 12rpx; }
.cat { padding: 14rpx; background: var(--eapp-bg); border-radius: 14rpx; display: flex; justify-content: space-between; align-items: center; font-size: 26rpx; }
.caret { color: var(--eapp-text-muted); }
.picker { padding: 14rpx; background: var(--eapp-bg); border-radius: 14rpx; font-size: 26rpx; }
.row { display: flex; align-items: center; justify-content: space-between; padding: 8rpx 0; font-size: 26rpx; }
.save { margin-top: 14rpx; }
</style>
