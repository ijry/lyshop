<script setup lang="ts">
import { reactive, watch } from 'vue'
const props = defineProps<{ show: boolean; deliveryMode?: 'express'|'local'|'both'; loading?: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'submit', payload: any): void }>()

const companyOptions = [
  { code: 'SF', name: '顺丰速运' }, { code: 'ZTO', name: '中通快递' }, { code: 'YTO', name: '圆通速递' },
  { code: 'STO', name: '申通快递' }, { code: 'YD', name: '韵达速递' }, { code: 'JD', name: '京东物流' },
  { code: 'EMS', name: '中国邮政 EMS' }, { code: 'DBL', name: '德邦快递' }, { code: 'JT', name: '极兔速递' },
]
const shipTypeOptions = [{ label: '首次发货', value: 'initial' }, { label: '补发', value: 'reship' }]
const deliveryTypeOptions = [{ label: '快递发货', value: 'express' }, { label: '同城配送', value: 'local' }]

const form = reactive({ ship_type: 'initial', delivery_type: 'express', company: 'SF', tracking_no: '', rider_name: '', rider_phone: '', after_sale_case_id: '', remark: '' })

watch(() => props.show, (v) => {
  if (v) {
    form.ship_type = 'initial'
    form.delivery_type = props.deliveryMode === 'both' ? 'express' : (props.deliveryMode || 'express')
    form.company = 'SF'; form.tracking_no = ''; form.rider_name = ''; form.rider_phone = ''
    form.after_sale_case_id = ''; form.remark = ''
  }
})

function onPick(field: keyof typeof form, options: any[], key: string, e: any) {
  ;(form as any)[field] = String(options[Number(e?.detail?.value || 0)]?.[key] || '')
}
function onSubmit() {
  if (form.ship_type === 'reship' && Number(form.after_sale_case_id || 0) <= 0) { uni.showToast({ title: '补发需填写售后单 ID', icon: 'none' }); return }
  if (form.delivery_type === 'local' && (!form.rider_name.trim() || !form.rider_phone.trim())) { uni.showToast({ title: '请填写骑手信息', icon: 'none' }); return }
  if (form.delivery_type === 'express' && (!form.company || !form.tracking_no.trim())) { uni.showToast({ title: '请填写快递与运单号', icon: 'none' }); return }
  emit('submit', {
    ship_type: form.ship_type, delivery_type: form.delivery_type,
    company: form.delivery_type === 'express' ? form.company : undefined,
    tracking_no: form.delivery_type === 'express' ? form.tracking_no.trim() : undefined,
    rider_name: form.delivery_type === 'local' ? form.rider_name.trim() : undefined,
    rider_phone: form.delivery_type === 'local' ? form.rider_phone.trim() : undefined,
    after_sale_case_id: form.ship_type === 'reship' ? Number(form.after_sale_case_id || 0) : undefined,
    remark: form.remark.trim() || undefined,
  })
}
</script>

<template>
  <up-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="popup">
      <view class="title">订单发货</view>
      <picker mode="selector" :range="shipTypeOptions" range-key="label" @change="(e) => onPick('ship_type', shipTypeOptions, 'value', e)">
        <view class="picker">{{ shipTypeOptions.find((x) => x.value === form.ship_type)?.label }}</view>
      </picker>
      <picker v-if="deliveryMode === 'both'" mode="selector" :range="deliveryTypeOptions" range-key="label" @change="(e) => onPick('delivery_type', deliveryTypeOptions, 'value', e)">
        <view class="picker mt">{{ deliveryTypeOptions.find((x) => x.value === form.delivery_type)?.label }}</view>
      </picker>
      <template v-if="form.delivery_type === 'express'">
        <picker mode="selector" :range="companyOptions" range-key="name" @change="(e) => onPick('company', companyOptions, 'code', e)">
          <view class="picker mt">{{ companyOptions.find((x) => x.code === form.company)?.name }}</view>
        </picker>
        <up-input v-model="form.tracking_no" placeholder="运单号" class="mt" />
      </template>
      <template v-else>
        <up-input v-model="form.rider_name" placeholder="骑手姓名" class="mt" />
        <up-input v-model="form.rider_phone" placeholder="骑手电话" class="mt" />
      </template>
      <up-input v-if="form.ship_type === 'reship'" v-model="form.after_sale_case_id" type="number" placeholder="售后单 ID" class="mt" />
      <up-input v-model="form.remark" placeholder="备注（可选）" class="mt" />
      <up-button type="primary" :loading="loading" class="mt-l" @click="onSubmit">确认发货</up-button>
    </view>
  </up-popup>
</template>

<style scoped>
.popup { padding: 24rpx; }
.title { font-size: 30rpx; font-weight: 700; margin-bottom: 14rpx; }
.picker { min-height: 76rpx; border: 1px solid var(--eapp-border); border-radius: 12rpx; padding: 0 20rpx; display: flex; align-items: center; }
.mt { margin-top: 12rpx; }
.mt-l { margin-top: 20rpx; }
</style>
