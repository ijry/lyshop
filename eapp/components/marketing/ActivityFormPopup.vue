<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import AppIcon from '@/components/common/AppIcon.vue'
import FormSheet from './FormSheet.vue'
import { activityMeta, formatActivityTime, type MarketingKind } from '@/utils/marketing'

const props = withDefaults(
  defineProps<{
    show: boolean
    kind: MarketingKind
    activity?: any
    submitting?: boolean
  }>(),
  { submitting: false },
)

const emit = defineEmits<{ close: []; submit: [payload: any] }>()

const meta = activityMeta(props.kind)
const form = reactive({
  name: '',
  start_at: '',
  end_at: '',
  status: 1,
})
const showStartPicker = ref(false)
const showEndPicker = ref(false)

watch(
  () => [props.show, props.activity],
  () => {
    if (!props.show) return
    form.name = String(props.activity?.name || '')
    form.start_at = String(props.activity?.start_at || '')
    form.end_at = String(props.activity?.end_at || '')
    form.status = props.activity ? (Number(props.activity?.status || 0) === 1 ? 1 : 0) : 1
  },
  { immediate: true },
)

function toISO(value: any): string {
  return new Date(value).toISOString()
}

function onSubmit() {
  const name = form.name.trim()
  if (!name) {
    uni.showToast({ title: '请输入活动名称', icon: 'none' })
    return
  }
  if (!form.start_at || !form.end_at) {
    uni.showToast({ title: '请选择活动起止时间', icon: 'none' })
    return
  }
  const start = new Date(form.start_at).getTime()
  const end = new Date(form.end_at).getTime()
  if (!start || !end || start >= end) {
    uni.showToast({ title: '结束时间需晚于开始时间', icon: 'none' })
    return
  }
  emit('submit', {
    name,
    start_at: form.start_at,
    end_at: form.end_at,
    status: Number(form.status),
  })
}
</script>

<template>
  <FormSheet
    :show="show"
    :title="activity ? `编辑${meta.title}` : `新建${meta.title}`"
    :subtitle="meta.desc"
    :submitting="submitting"
    submit-text="保存活动"
    body-height="44vh"
    @close="emit('close')"
    @submit="onSubmit"
  >
    <view class="eapp-form-group">
      <view class="eapp-form-group__title">
        <view class="eapp-form-group__bar" />
        <text>基础信息</text>
      </view>

      <text class="eapp-form-label">
        活动名称<text class="eapp-form-label__required">*</text>
      </text>
      <view class="eapp-form-input">
        <up-input v-model="form.name" placeholder="例如：周末限时秒杀" clearable :maxlength="30" />
      </view>

      <text class="eapp-form-label">活动状态</text>
      <view class="eapp-form-segment">
        <view
          class="eapp-form-segment__item"
          :class="{ 'eapp-form-segment__item--active': form.status === 1 }"
          hover-class="eapp-hover"
          @click="form.status = 1"
        >
          启用
        </view>
        <view
          class="eapp-form-segment__item"
          :class="{ 'eapp-form-segment__item--active': form.status === 0 }"
          hover-class="eapp-hover"
          @click="form.status = 0"
        >
          停用
        </view>
      </view>
    </view>

    <view class="eapp-form-group">
      <view class="eapp-form-group__title">
        <view class="eapp-form-group__bar" />
        <text>活动周期</text>
      </view>

      <text class="eapp-form-label">
        开始时间<text class="eapp-form-label__required">*</text>
      </text>
      <view class="eapp-form-picker" hover-class="eapp-hover" @click="showStartPicker = true">
        <text class="eapp-form-picker__value" :class="{ 'eapp-form-picker__value--empty': !form.start_at }">
          {{ form.start_at ? formatActivityTime(form.start_at, true) : '请选择开始时间' }}
        </text>
        <AppIcon name="chevron-right" :size="26" color="faint" />
      </view>

      <text class="eapp-form-label">
        结束时间<text class="eapp-form-label__required">*</text>
      </text>
      <view class="eapp-form-picker" hover-class="eapp-hover" @click="showEndPicker = true">
        <text class="eapp-form-picker__value" :class="{ 'eapp-form-picker__value--empty': !form.end_at }">
          {{ form.end_at ? formatActivityTime(form.end_at, true) : '请选择结束时间' }}
        </text>
        <AppIcon name="chevron-right" :size="26" color="faint" />
      </view>

      <view class="eapp-form-tip">
        <AppIcon name="info" :size="24" color="primary" />
        <text>同一类型活动的时间段不可重叠，保存时会自动校验。</text>
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
</template>
