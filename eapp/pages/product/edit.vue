<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import SkuMatrixEditor from '@/components/biz/SkuMatrixEditor.vue'
import CategoryTreePicker from '@/components/biz/CategoryTreePicker.vue'
import RichTextEditor from '@/components/biz/RichTextEditor.vue'
import { createProduct, getProductDetail, updateProduct } from '@/api/product'
import { getSpecTemplates } from '@/api/spec-template'

const id = ref(0)
const saving = ref(false)
const showCatPicker = ref(false)
const showOnlinePicker = ref(false)
const showOfflinePicker = ref(false)
const showShippingPicker = ref(false)
const shippingOptions = [
  { label: '默认模板', value: 'default' },
  { label: '包邮', value: 'free' },
  { label: '到付', value: 'cod' },
  { label: '同城', value: 'local' },
]

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

const specTemplates = ref<any[]>([])
const showSpecPicker = ref(false)

type SkuAttr = { name: string; value: string }

function normalizeSkuAttrs(attrs: any): SkuAttr[] {
  const list = Array.isArray(attrs) ? attrs : []
  return list
    .map((item: any) => ({
      name: String(item?.name || '').trim(),
      value: String(item?.value || '').trim(),
    }))
    .filter((item: SkuAttr) => item.name && item.value)
    .sort((a: SkuAttr, b: SkuAttr) => {
      if (a.name === b.name) return a.value.localeCompare(b.value)
      return a.name.localeCompare(b.name)
    })
}

function buildSkuKey(attrs: SkuAttr[]) {
  const normalized = normalizeSkuAttrs(attrs)
  if (!normalized.length) return '__default__'
  return normalized.map((item: SkuAttr) => `${item.name}:${item.value}`).join('|')
}

function buildSpecSchema(rows: any[]) {
  const groupMap = new Map<string, Set<string>>()
  for (const row of rows) {
    const attrs = normalizeSkuAttrs(row?.attrs)
    for (const attr of attrs) {
      if (!groupMap.has(attr.name)) groupMap.set(attr.name, new Set<string>())
      groupMap.get(attr.name)?.add(attr.value)
    }
  }
  return Array.from(groupMap.entries())
    .map(([name, values]) => ({
      name,
      values: Array.from(values).sort((a, b) => a.localeCompare(b)),
    }))
    .sort((a, b) => a.name.localeCompare(b.name))
}

function buildSkuOverrides(rows: any[]) {
  const dedup = new Map<string, { sku_key: string; sku_code: string; price: number; stock: number }>()
  for (const row of rows) {
    const attrs = normalizeSkuAttrs(row?.attrs)
    const skuKey = buildSkuKey(attrs)
    dedup.set(skuKey, {
      sku_key: skuKey,
      sku_code: String(row?.sku_code || '').trim(),
      price: Number(row?.price || 0),
      stock: Number(row?.stock || 0),
    })
  }
  return Array.from(dedup.values())
}

function buildDetailPayload() {
  const html = String(form.detail_html || '').trim()
  if (!html) return { version: 1, blocks: [] as Array<Record<string, any>> }
  return {
    version: 1,
    blocks: [{ id: `detail-${Date.now()}`, type: 'rich_text', props: { content: html } }],
  }
}

function parseDetailHTML(detail: any): string {
  if (!detail) return ''
  const payload = typeof detail === 'string'
    ? (() => { try { return JSON.parse(detail) } catch { return null } })()
    : detail
  if (!payload || !Array.isArray(payload.blocks)) return ''
  const richText = payload.blocks.find((block: any) => String(block?.type || '') === 'rich_text')
  if (richText?.props?.content) return String(richText.props.content)
  const textBlock = payload.blocks.find((block: any) => String(block?.type || '') === 'text')
  if (textBlock?.props?.text) return String(textBlock.props.text)
  return ''
}

function mapImagesToCovers(data: any): string[] {
  const images = Array.isArray(data?.images)
    ? data.images.map((item: any) => String(item?.url || '').trim()).filter((url: string) => !!url)
    : []
  if (images.length) return images
  if (Array.isArray(data?.covers)) {
    return data.covers.map((item: any) => String(item || '').trim()).filter((url: string) => !!url)
  }
  if (String(data?.cover || '').trim()) return [String(data.cover).trim()]
  return []
}

async function loadSpecTemplates() {
  const data: any = await getSpecTemplates({ page: 1, size: 200 })
  specTemplates.value = Array.isArray(data?.list) ? data.list : []
}

function applySpecTemplate(tpl: any) {
  if (!tpl?.attrs || !Array.isArray(tpl.attrs)) return
  const skus: Array<{ attrs: Array<{ name: string; value: string }>; price: number; stock: number }> = []
  function cross(idx: number): Array<Array<{ name: string; value: string }>> {
    if (idx >= tpl.attrs.length) return [[]]
    const sub = cross(idx + 1)
    return (tpl.attrs[idx].values || []).flatMap((v: string) => sub.map((row: any) => [{ name: tpl.attrs[idx].name, value: v }, ...row]))
  }
  const combos = cross(0)
  for (const attrs of combos) {
    skus.push({ attrs, price: Number(form.price || 0), stock: 0 })
  }
  form.skus = skus
  showSpecPicker.value = false
  uni.showToast({ title: '已应用模板', icon: 'success' })
}

async function loadData() {
  if (!id.value) return
  const data: any = await getProductDetail(id.value)
  Object.assign(form, {
    title: data?.title || '', subtitle: data?.subtitle || '',
    sell_points: Array.isArray(data?.sell_points) ? data.sell_points : [],
    covers: mapImagesToCovers(data),
    detail_html: parseDetailHTML(data?.detail),
    price: Number(data?.price || 0), stock: Number(data?.stock || 0),
    unit: String(data?.unit || '件'), weight: Number(data?.weight || 0),
    category_id: Number(data?.category_id || 0),
    category_path_name: String(data?.category_path_name || ''),
    tags: Array.isArray(data?.tags) ? data.tags : [],
    skus: Array.isArray(data?.skus)
      ? data.skus.map((sku: any) => ({
        ...sku,
        attrs: normalizeSkuAttrs(sku?.attrs),
        price: Number(sku?.price || 0),
        stock: Number(sku?.stock || 0),
      }))
      : [],
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
    const specSchema = buildSpecSchema(form.skus)
    const skuOverrides = buildSkuOverrides(form.skus)
    const payload: any = {
      product: {
        title: form.title.trim(),
        subtitle: form.subtitle.trim(),
        cover: form.covers[0] || '',
        price: Number(form.price),
        stock: Number(form.stock),
        category_id: Number(form.category_id),
        status: Number(form.status || 0) === 1 ? 1 : 0,
        detail: buildDetailPayload(),
      },
      images: form.covers
        .map((url: string, idx: number) => ({ url: String(url || '').trim(), sort: idx }))
        .filter((item: any) => !!item.url),
      sku_generation_mode: 'auto',
      spec_schema: specSchema,
      sku_overrides: skuOverrides,
    }
    const res: any = id.value ? await updateProduct(id.value, payload) : await createProduct(payload)
    const diff = res?.sku_diff
    if (diff && typeof diff === 'object') {
      uni.showToast({
        title: `保存成功 +${Number(diff.added || 0)} / =${Number(diff.kept || 0)} / -${Number(diff.inactivated || 0)}`,
        icon: 'none',
        duration: 2500,
      })
    } else {
      uni.showToast({ title: '保存成功', icon: 'success' })
    }
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
      <view class="section-title">规格 SKU
        <up-button size="mini" plain @click="loadSpecTemplates().then(() => showSpecPicker = true)">应用规格模板</up-button>
      </view>
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
      <view class="picker" @click="showShippingPicker = true">物流模板：{{ { default: '默认模板', free: '包邮', cod: '到付', local: '同城' }[form.shipping_template] }}</view>
      <up-picker :show="showShippingPicker" :columns="[shippingOptions]" keyName="label" @confirm="(e) => { form.shipping_template = e.value[0].value; showShippingPicker = false }" @cancel="showShippingPicker = false" @close="showShippingPicker = false" />
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
      <view class="picker" @click="showOnlinePicker = true">{{ form.online_at ? form.online_at.slice(0, 16).replace('T', ' ') : '上架时间（可选）' }}</view>
      <up-datetime-picker :show="showOnlinePicker" v-model="form.online_at" mode="datetime" @confirm="(e) => { form.online_at = new Date(e.value).toISOString(); showOnlinePicker = false }" @cancel="showOnlinePicker = false" @close="showOnlinePicker = false" />
      <view class="picker mt" @click="showOfflinePicker = true">{{ form.offline_at ? form.offline_at.slice(0, 16).replace('T', ' ') : '下架时间（可选）' }}</view>
      <up-datetime-picker :show="showOfflinePicker" v-model="form.offline_at" mode="datetime" @confirm="(e) => { form.offline_at = new Date(e.value).toISOString(); showOfflinePicker = false }" @cancel="showOfflinePicker = false" @close="showOfflinePicker = false" />
    </view>

    <up-button type="primary" :loading="saving" class="save" @click="save">保存</up-button>

    <CategoryTreePicker :show="showCatPicker" :value="form.category_id" @close="showCatPicker = false" @pick="onPickCategory" />

    <up-popup :show="showSpecPicker" mode="bottom" round="16" @close="showSpecPicker = false">
      <view class="popup-body">
        <view class="popup-title">选择规格模板</view>
        <view v-if="!specTemplates.length" class="empty-tpl">暂无模板</view>
        <view v-for="tpl in specTemplates" :key="tpl.id" class="tpl-item" @click="applySpecTemplate(tpl)">
          <text class="tpl-name">{{ tpl.name }}</text>
          <text class="tpl-desc">{{ (tpl.attrs || []).map((a: any) => a.name).join(' / ') }}</text>
        </view>
      </view>
    </up-popup>
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
.popup-body { padding: 24rpx; box-sizing: border-box; }
.popup-title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.empty-tpl { text-align: center; color: var(--eapp-text-muted); padding: 40rpx 0; }
.tpl-item { padding: 16rpx; border: 1px solid var(--eapp-border); border-radius: 14rpx; margin-bottom: 10rpx; }
.tpl-name { font-size: 28rpx; font-weight: 600; }
.tpl-desc { display: block; margin-top: 6rpx; color: var(--eapp-text-muted); font-size: 22rpx; }
</style>
