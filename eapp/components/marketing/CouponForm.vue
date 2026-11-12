<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import AppIcon from '@/components/common/AppIcon.vue'
import FormSheet from './FormSheet.vue'
import { formatActivityTime } from '@/utils/marketing'

const props = withDefaults(
  defineProps<{ show: boolean; coupon?: any; submitting?: boolean }>(),
  { submitting: false },
)
const emit = defineEmits<{ close: []; submit: [payload: any] }>()

const form = reactive({
  name: '',
  type: 1,
  min_amount: '',
  discount: '',
  total_count: '',
  per_limit: '1',
  status: 1,
  start_at: '',
  end_at: '',
  description: '',
  stack_rule: 'exclusive',
  target_type: 'all',
  target_value: '',
})

const TYPE_OPTIONS = [
  { label: '满减', value: 1 },
  { label: '折扣', value: 2 },
  { label: '立减', value: 3 },
]
const STACK_OPTIONS = [
  { label: '互斥（不与其他券叠加）', value: 'exclusive' },
  { label: '同类可叠', value: 'same_type' },
  { label: '跨类可叠', value: 'cross_type' },
]
const TARGET_OPTIONS = [
  { label: '全部用户', value: 'all' },
  { label: '会员等级', value: 'vip_level' },
  { label: '新用户', value: 'new_user' },
]

const showStartPicker = ref(false)
const showEndPicker = ref(false)
const showStackPicker = ref(false)
const showTargetPicker = ref(false)

watch(
  () => [props.show, props.coupon],
  () => {
    if (!props.show) return
    const c = props.coupon
    form.name = String(c?.name || '')
    form.type = Number(c?.type || 1)
    form.min_amount = String(c?.min_amount ?? '')
    form.discount = String(c?.discount ?? '')
    form.total_count = String(c?.total_count ?? '')
    form.per_limit = String(c?.per_limit ?? '1')
    form.status = c ? (Number(c?.status || 0) === 1 ? 1 : 0) : 1
    form.start_at = String(c?.start_at || '')
    form.end_at = String(c?.end_at || '')
    form.description = String(c?.description || '')
    form.stack_rule = String(c?.stack_rule || 'exclusive')
    form.target_type = String(c?.target_type || 'all')
    form.target_value = String(c?.target_value || '')
  },
  { immediate: true },
)

function optionLabel(options: Array<{ label: string; value: any }>, value: any) {
  return options.find((item) => item.value === value)?.label || '请选择'
}

function toISO(value: any): string {
  return new Date(value).toISOString()
}

function onSubmit() {
  const name = form.name.trim()
  if (!name) {
    uni.showToast({ title: '请输入优惠券名称', icon: 'none' })
    return
  }
  if (Number(form.discount || 0) <= 0) {
    uni.showToast({ title: '请输入有效的优惠值', icon: 'none' })
    return
  }
  if (form.start_at && form.end_at) {
    const start = new Date(form.start_at).getTime()
    const end = new Date(form.end_at).getTime()
    if (!start || !end || start >= end) {
      uni.showToast({ title: '结束时间需晚于开始时间', icon: 'none' })
      return
    }
  }
  emit('submit', {
    name,
    type: Number(form.type),
    min_amount: Number(form.min_amount || 0),
    discount: Number(form.discount || 0),
    total_count: Number(form.total_count || 0),
    per_limit: Number(form.per_limit || 1),
    status: Number(form.status),
    start_at: form.start_at,
    end_at: form.end_at,
    description: form.description.trim(),
    stack_rule: form.stack_rule,
    target_type: form.target_type,
    target_value: form.target_value,
  })
}
</script>

<template>
  <FormSheet
    :show="show"
    :title="coupon ? '编辑优惠券' : '新建优惠券'"
    subtitle="配置优惠规则、发放范围与有效期"
    :submitting="submitting"
    submit-text="保存优惠券"
    body-height="64vh"
    @close="emit('close')"
    @submit="onSubmit"
  >
    <view class="eapp-form-group">
      <view class="eapp-form-group__title">
        <view class="eapp-form-group__bar" />
        <text>基础信息</text>
      </view>

      <text class="eapp-form-label">券名称<text class="eapp-form-label__required">*</text></text>
      <view class="eapp-form-input">
        <up-input v-model="form.name" placeholder="例如：新人满减券" clearable :maxlength="30" />
      </view>

      <text class="eapp-form-label">优惠类型</text>
      <view class="eapp-form-segment">
        <view
          v-for="item in TYPE_OPTIONS"
          :key="item.value"
          class="eapp-form-segment__item"
          :class="{ 'eapp-form-segment__item--active': form.type === item.value }"
          hover-class="eapp-hover"
          @click="form.type = item.value"
        >
          {{ item.label }}
        </view>
      </view>

      <text class="eapp-form-label">状态</text>
      <view class="eapp-form-segment">
        <view
          class="eapp-form-segment__item"
          :class="{ 'eapp-form-segment__item--active': form.status === 1 }"
          hover-class="eapp-hover"
          @click="form.status = 1"
        >启用</view>
        <view
          class="eapp-form-segment__item"
          :class="{ 'eapp-form-segment__item--active': form.status === 0 }"
          hover-class="eapp-hover"
          @click="form.status = 0"
        >停用</view>
      </view>
    </view>

    <view class="eapp-form-group">
      <view class="eapp-form-group__title">
        <view class="eapp-form-group__bar" />
        <text>优惠规则</text>
      </view>

      <view class="eapp-form-row">
        <view class="eapp-form-col">
          <text class="eapp-form-label">优惠值<text class="eapp-form-label__required">*</text></text>
          <view class="eapp-form-input">
            <up-input v-model="form.discount" type="digit" placeholder="金额 / 折扣" />
          </view>
        </view>
        <view class="eapp-form-col">
          <text class="eapp-form-label">最低消费（元）</text>
          <view class="eapp-form-input">
            <up-input v-model="form.min_amount" type="digit" placeholder="0 为不限" />
          </view>
        </view>
      </view>
      <text class="eapp-form-hint">满减 / 立减填金额，折扣填折扣数（如 8.5 表示 8.5 折）。</text>

      <text class="eapp-form-label">叠加规则</text>
      <view class="eapp-form-picker" hover-class="eapp-hover" @click="showStackPicker = true">
        <text class="eapp-form-picker__value">{{ optionLabel(STACK_OPTIONS, form.stack_rule) }}</text>
        <AppIcon name="chevron-right" :size="26" color="faint" />
      </view>
    </view>

    <view class="eapp-form-group">
      <view class="eapp-form-group__title">
        <view class="eapp-form-group__bar" />
        <text>发放设置</text>
      </view>

      <view class="eapp-form-row">
        <view class="eapp-form-col">
          <text class="eapp-form-label">发行总量</text>
          <view class="eapp-form-input">
            <up-input v-model="form.total_count" type="number" placeholder="0 为不限" />
          </view>
        </view>
        <view class="eapp-form-col">
          <text class="eapp-form-label">每人限领</text>
          <view class="eapp-form-input">
            <up-input v-model="form.per_limit" type="number" placeholder="1" />
          </view>
        </view>
      </view>

      <text class="eapp-form-label">生效时间</text>
      <view class="eapp-form-picker" hover-class="eapp-hover" @click="showStartPicker = true">
        <text class="eapp-form-picker__value" :class="{ 'eapp-form-picker__value--empty': !form.start_at }">
          {{ form.start_at ? formatActivityTime(form.start_at, true) : '不限制开始时间' }}
        </text>
        <AppIcon name="chevron-right" :size="26" color="faint" />
      </view>

      <text class="eapp-form-label">失效时间</text>
      <view class="eapp-form-picker" hover-class="eapp-hover" @click="showEndPicker = true">
        <text class="eapp-form-picker__value" :class="{ 'eapp-form-picker__value--empty': !form.end_at }">
          {{ form.end_at ? formatActivityTime(form.end_at, true) : '不限制结束时间' }}
        </text>
        <AppIcon name="chevron-right" :size="26" color="faint" />
      </view>

      <text class="eapp-form-label">目标用户</text>
      <view class="eapp-form-picker" hover-class="eapp-hover" @click="showTargetPicker = true">
        <text class="eapp-form-picker__value">{{ optionLabel(TARGET_OPTIONS, form.target_type) }}</text>
        <AppIcon name="chevron-right" :size="26" color="faint" />
      </view>

      <template v-if="form.target_type === 'vip_level'">
        <text class="eapp-form-label">会员等级值</text>
        <view class="eapp-form-input">
          <up-input v-model="form.target_value" placeholder="例如：3" />
        </view>
      </template>

      <text class="eapp-form-label">说明</text>
      <view class="eapp-form-input">
        <up-input v-model="form.description" placeholder="展示给用户的券说明" />
      </view>
    </view>
  </FormSheet>

  <up-datetime-picker
    :show="showStartPicker"
    v-model="form.start_at"
    mode="datetime"
    @confirm="(e: any) => { form.start_at = toISO(e.value); showStartPicker = false }"
    @cancel="showStartPicker = false"
    @close="showStartPicker = false"
  />
  <up-datetime-picker
    :show="showEndPicker"
    v-model="form.end_at"
    mode="datetime"
    @confirm="(e: any) => { form.end_at = toISO(e.value); showEndPicker = false }"
    @cancel="showEndPicker = false"
    @close="showEndPicker = false"
  />
  <up-picker
    :show="showStackPicker"
    :columns="[STACK_OPTIONS]"
    keyName="label"
    @confirm="(e: any) => { form.stack_rule = String(e?.value?.[0]?.value || 'exclusive'); showStackPicker = false }"
    @cancel="showStackPicker = false"
    @close="showStackPicker = false"
  />
  <up-picker
    :show="showTargetPicker"
    :columns="[TARGET_OPTIONS]"
    keyName="label"
    @confirm="(e: any) => { form.target_type = String(e?.value?.[0]?.value || 'all'); showTargetPicker = false }"
    @cancel="showTargetPicker = false"
    @close="showTargetPicker = false"
  />
</template>
