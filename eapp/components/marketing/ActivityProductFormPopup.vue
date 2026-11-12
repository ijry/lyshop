<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import AppIcon from '@/components/common/AppIcon.vue'
import FormSheet from './FormSheet.vue'
import { activityMeta, formatMoney, type MarketingKind } from '@/utils/marketing'
import type { SkuOption } from '@/composables/useActivityCenter'

const props = withDefaults(
  defineProps<{
    show: boolean
    kind: MarketingKind
    row?: any
    products?: any[]
    skus?: SkuOption[]
    submitting?: boolean
  }>(),
  { products: () => [], skus: () => [], submitting: false },
)

const emit = defineEmits<{
  close: []
  submit: [payload: any]
  'select-product': [productID: number]
}>()

const meta = computed(() => activityMeta(props.kind))
const form = reactive({
  product_id: 0,
  sku_id: 0,
  activity_price: 0,
  start_price: 0,
  floor_price: 0,
  limit_per_order: 0,
  total_stock_limit: 0,
})

const showProductPicker = ref(false)
const showSkuPicker = ref(false)

const productList = computed(() =>
  props.products.length ? props.products : [{ id: 0, title: '暂无可选商品' }],
)
const skuList = computed(() => (props.skus.length ? props.skus : [{ id: 0, label: '请先选择商品' }]))
const productName = computed(() => {
  const hit = props.products.find((item) => Number(item.id || 0) === Number(form.product_id || 0))
  return hit ? String(hit.title || `商品 #${hit.id}`) : ''
})
const skuLabel = computed(() => {
  const hit = props.skus.find((item) => Number(item.id || 0) === Number(form.sku_id || 0))
  return hit ? String(hit.label) : ''
})
const originPrice = computed(() => {
  const hit = props.skus.find((item) => Number(item.id || 0) === Number(form.sku_id || 0))
  return Number(hit?.price || 0)
})

watch(
  () => [props.show, props.row],
  () => {
    if (!props.show) return
    form.product_id = Number(props.row?.product_id || 0)
    form.sku_id = Number(props.row?.sku_id || 0)
    form.activity_price = Number(props.row?.activity_price || 0)
    form.start_price = Number(props.row?.start_price || 0)
    form.floor_price = Number(props.row?.floor_price || 0)
    form.limit_per_order = Number(props.row?.limit_per_order || 0)
    form.total_stock_limit = Number(props.row?.total_stock_limit || 0)
  },
  { immediate: true },
)

function pickProduct(event: any) {
  const id = Number(event?.value?.[0]?.id || 0)
  form.product_id = id
  form.sku_id = 0
  showProductPicker.value = false
  if (id) emit('select-product', id)
}

function onSubmit() {
  if (!form.product_id || !form.sku_id) {
    uni.showToast({ title: '请选择商品与 SKU', icon: 'none' })
    return
  }
  if (props.kind === 'bargain') {
    if (form.start_price <= 0 || form.floor_price <= 0 || form.floor_price > form.start_price) {
      uni.showToast({ title: '请检查起砍价与底价', icon: 'none' })
      return
    }
  } else if (form.activity_price <= 0) {
    uni.showToast({ title: '请输入有效的活动价', icon: 'none' })
    return
  }
  if (form.limit_per_order < 0 || form.total_stock_limit < 0) {
    uni.showToast({ title: '限购与库存不能为负', icon: 'none' })
    return
  }
  emit('submit', {
    product_id: form.product_id,
    sku_id: form.sku_id,
    activity_price: Number(form.activity_price),
    start_price: Number(form.start_price),
    floor_price: Number(form.floor_price),
    limit_per_order: Number(form.limit_per_order),
    total_stock_limit: Number(form.total_stock_limit),
  })
}
</script>

<template>
  <FormSheet
    :show="show"
    :title="row ? '编辑活动商品' : '添加活动商品'"
    :subtitle="meta.title"
    :submitting="submitting"
    submit-text="保存商品"
    body-height="60vh"
    @close="emit('close')"
    @submit="onSubmit"
  >
    <view class="eapp-form-group">
      <view class="eapp-form-group__title">
        <view class="eapp-form-group__bar" />
        <text>商品与规格</text>
      </view>

      <text class="eapp-form-label">选择商品<text class="eapp-form-label__required">*</text></text>
      <view class="eapp-form-picker" hover-class="eapp-hover" @click="showProductPicker = true">
        <text class="eapp-form-picker__value" :class="{ 'eapp-form-picker__value--empty': !productName }">
          {{ productName || '请选择商品' }}
        </text>
        <AppIcon name="chevron-right" :size="26" color="faint" />
      </view>

      <text class="eapp-form-label">选择 SKU<text class="eapp-form-label__required">*</text></text>
      <view class="eapp-form-picker" hover-class="eapp-hover" @click="showSkuPicker = true">
        <text class="eapp-form-picker__value" :class="{ 'eapp-form-picker__value--empty': !skuLabel }">
          {{ skuLabel || '请选择 SKU' }}
        </text>
        <AppIcon name="chevron-right" :size="26" color="faint" />
      </view>
      <text v-if="originPrice" class="eapp-form-hint eapp-num">商品原价 {{ formatMoney(originPrice) }}</text>
    </view>

    <view class="eapp-form-group">
      <view class="eapp-form-group__title">
        <view class="eapp-form-group__bar" />
        <text>价格规则</text>
      </view>

      <template v-if="kind !== 'bargain'">
        <text class="eapp-form-label">{{ meta.priceLabel }}（元）<text class="eapp-form-label__required">*</text></text>
        <view class="eapp-form-input">
          <up-input v-model="form.activity_price" type="digit" placeholder="0.00" />
        </view>
      </template>
      <template v-else>
        <view class="eapp-form-row">
          <view class="eapp-form-col">
            <text class="eapp-form-label">起砍价（元）<text class="eapp-form-label__required">*</text></text>
            <view class="eapp-form-input">
              <up-input v-model="form.start_price" type="digit" placeholder="0.00" />
            </view>
          </view>
          <view class="eapp-form-col">
            <text class="eapp-form-label">底价（元）<text class="eapp-form-label__required">*</text></text>
            <view class="eapp-form-input">
              <up-input v-model="form.floor_price" type="digit" placeholder="0.00" />
            </view>
          </view>
        </view>
      </template>
      <text class="eapp-form-hint">{{ meta.priceHint }}</text>
    </view>

    <view class="eapp-form-group">
      <view class="eapp-form-group__title">
        <view class="eapp-form-group__bar" />
        <text>销售限制</text>
      </view>
      <view class="eapp-form-row">
        <view class="eapp-form-col">
          <text class="eapp-form-label">每单限购</text>
          <view class="eapp-form-input">
            <up-input v-model="form.limit_per_order" type="number" placeholder="0" />
          </view>
        </view>
        <view class="eapp-form-col">
          <text class="eapp-form-label">活动库存</text>
          <view class="eapp-form-input">
            <up-input v-model="form.total_stock_limit" type="number" placeholder="0" />
          </view>
        </view>
      </view>
      <view class="eapp-form-tip">
        <AppIcon name="info" :size="24" color="primary" />
        <text>填 0 表示不限制；同一活动下同一 SKU 只能出现一次。</text>
      </view>
    </view>
  </FormSheet>

  <up-picker
    :show="showProductPicker"
    :columns="[productList]"
    keyName="title"
    @confirm="pickProduct"
    @cancel="showProductPicker = false"
    @close="showProductPicker = false"
  />
  <up-picker
    :show="showSkuPicker"
    :columns="[skuList]"
    keyName="label"
    @confirm="(e: any) => { form.sku_id = Number(e?.value?.[0]?.id || 0); showSkuPicker = false }"
    @cancel="showSkuPicker = false"
    @close="showSkuPicker = false"
  />
</template>
