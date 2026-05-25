<template>
  <view class="min-h-screen" style="background: #f5f5f5; padding: 12px 16px;">
    <view style="background: #fff; border-radius: 16px; overflow: hidden; box-shadow: 0 2px 12px rgba(0,0,0,0.04);">
      <view @click="showDeleteConfirm = true"
        style="display: flex; align-items: center; padding: 14px 20px;">
        <u-icon name="lock" size="20" color="#dc2626" />
        <text style="flex: 1; margin-left: 12px; font-size: 14px; color: #dc2626;">{{ $t('security.deleteAccount') }}</text>
        <u-icon name="arrow-right" size="14" color="#ccc" />
      </view>
    </view>

    <view style="padding: 12px 4px;">
      <text style="font-size: 12px; color: #999; line-height: 1.6;">
        {{ $t('security.deleteWarning') }}
      </text>
    </view>

    <u-popup :show="showDeleteConfirm" mode="center" round="20" @close="showDeleteConfirm = false">
      <view style="padding: 30px; width: 300px;">
        <text style="font-size: 17px; font-weight: 700; color: #111; display: block; text-align: center;">{{ $t('security.deleteAccount') }}</text>
        <text style="font-size: 13px; color: #999; display: block; text-align: center; margin: 12px 0 20px; line-height: 1.5;">
          {{ $t('security.deleteConfirmMsg') }}
        </text>
        <view style="margin-bottom: 12px;">
          <u-input v-model="deleteForm.phone" :placeholder="$t('security.phone')" type="number" :maxlength="11" border="surround" shape="circle" />
        </view>
        <view style="display: flex; gap: 10px; margin-bottom: 20px;">
          <view style="flex: 1;">
            <u-input v-model="deleteForm.code" :placeholder="$t('security.verifyCode')" type="number" :maxlength="6" border="surround" shape="circle" />
          </view>
          <u-button size="small" :disabled="deleteCountdown > 0"
            :text="deleteCountdown > 0 ? `${deleteCountdown}s` : $t('security.getCode')"
            @click="sendDeleteCode" type="primary" plain shape="circle" />
        </view>
        <view style="display: flex; gap: 10px;">
          <u-button :text="$t('security.cancel')" @click="showDeleteConfirm=false" shape="circle" class="flex-1" />
          <u-button :text="$t('security.confirmDelete')" type="error" @click="deleteAccount" shape="circle" class="flex-1"
            :custom-style="{background: '#dc2626', borderColor: '#dc2626'}" />
        </view>
      </view>
    </u-popup>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { post } from '@/utils/request'

const { t } = useI18n()

const showDeleteConfirm = ref(false)
const deleteForm = ref({ phone: '', code: '' })
const deleteCountdown = ref(0)

async function sendDeleteCode() {
  if (!deleteForm.value.phone || deleteForm.value.phone.length !== 11) {
    uni.showToast({ title: t('security.phoneRequired'), icon: 'none' })
    return
  }
  try {
    const data = await post<any>('/api/v1/auth/sms/send', { phone: deleteForm.value.phone })
    if (data?.dev_code) deleteForm.value.code = data.dev_code
  } catch {}
  deleteCountdown.value = 60
  const timer = setInterval(() => {
    if (--deleteCountdown.value <= 0) clearInterval(timer)
  }, 1000)
}

async function deleteAccount() {
  if (!deleteForm.value.code) {
    uni.showToast({ title: t('security.codeRequired'), icon: 'none' })
    return
  }
  try {
    await post('/api/v1/user/delete', deleteForm.value)
    uni.showToast({ title: t('security.accountDeleted'), icon: 'success' })
    setTimeout(() => {
      uni.removeStorageSync('user_token')
      uni.reLaunch({ url: '/pages/login/index' })
    }, 1500)
  } catch (e: any) {
    uni.showToast({ title: e.message || t('security.deleteFailed'), icon: 'none' })
  }
}
</script>
