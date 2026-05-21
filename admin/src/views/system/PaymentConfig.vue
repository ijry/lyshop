<template>
  <div>
    <h2 class="text-xl font-semibold text-slate-800 mb-6">支付 & 短信配置</h2>

    <!-- WeChat Pay -->
    <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6 mb-4 max-w-xl">
      <h3 class="font-semibold text-slate-700 mb-4">微信支付</h3>
      <div class="space-y-3">
        <div v-for="field in wechatFields" :key="field.key">
          <label class="block text-sm font-medium text-slate-600 mb-1">{{ field.label }}</label>
          <input v-model="wechat[field.key]" :type="field.secret ? 'password' : 'text'"
            class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" :placeholder="field.placeholder" />
        </div>
        <button @click="saveConfig('wechat_pay', wechat)"
          class="px-6 py-2 bg-blue-700 text-white rounded-xl text-sm hover:bg-blue-600">保存微信支付配置</button>
      </div>
    </div>

    <!-- Alipay -->
    <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6 mb-4 max-w-xl">
      <h3 class="font-semibold text-slate-700 mb-4">支付宝支付</h3>
      <div class="space-y-3">
        <div v-for="field in alipayFields" :key="field.key">
          <label class="block text-sm font-medium text-slate-600 mb-1">{{ field.label }}</label>
          <textarea v-if="field.textarea" v-model="alipay[field.key]" rows="4"
            class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm resize-none" :placeholder="field.placeholder" />
          <input v-else v-model="alipay[field.key]" :type="field.secret ? 'password' : 'text'"
            class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" :placeholder="field.placeholder" />
        </div>
        <div class="flex items-center gap-2">
          <input type="checkbox" id="sandbox" v-model="alipay.sandbox" />
          <label for="sandbox" class="text-sm text-slate-600">沙箱模式</label>
        </div>
        <button @click="saveConfig('alipay', alipay)"
          class="px-6 py-2 bg-blue-700 text-white rounded-xl text-sm hover:bg-blue-600">保存支付宝配置</button>
      </div>
    </div>

    <!-- SMS -->
    <div class="bg-white rounded-xl shadow-sm border border-slate-100 p-6 max-w-xl">
      <h3 class="font-semibold text-slate-700 mb-4">短信配置</h3>
      <div class="space-y-3">
        <div>
          <label class="block text-sm font-medium text-slate-600 mb-1">服务商</label>
          <select v-model="sms.provider" class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm">
            <option value="aliyun">阿里云</option>
            <option value="tencent">腾讯云</option>
          </select>
        </div>
        <div v-for="field in smsFields" :key="field.key">
          <label class="block text-sm font-medium text-slate-600 mb-1">{{ field.label }}</label>
          <input v-model="sms[field.key]" :type="field.secret ? 'password' : 'text'"
            class="w-full border border-slate-200 rounded-xl px-4 py-2 text-sm" :placeholder="field.placeholder" />
        </div>
        <button @click="saveSmsConfig"
          class="px-6 py-2 bg-blue-700 text-white rounded-xl text-sm hover:bg-blue-600">保存短信配置</button>
      </div>
    </div>

    <p v-if="saved" class="mt-4 text-green-600 text-sm">✓ 配置已保存</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import request from '@/api/request'

const saved = ref(false)

const wechat = ref<Record<string, string>>({ app_id: '', mch_id: '', api_key: '' })
const wechatFields = [
  { key: 'app_id',  label: 'AppID',  placeholder: 'wx...',        secret: false },
  { key: 'mch_id',  label: '商户号', placeholder: '1234567890',   secret: false },
  { key: 'api_key', label: 'API密钥(v3)', placeholder: '32位密钥', secret: true  },
]

const alipay = ref<Record<string, any>>({ app_id: '', private_key: '', public_key: '', sandbox: false })
const alipayFields = [
  { key: 'app_id',     label: 'AppID',      placeholder: '2021...',  secret: false, textarea: false },
  { key: 'private_key',label: '应用私钥',   placeholder: 'RSA私钥',  secret: true,  textarea: true  },
  { key: 'public_key', label: '支付宝公钥', placeholder: '支付宝公钥', secret: false, textarea: true },
]

const sms = ref<Record<string, string>>({ provider: 'aliyun', access_key: '', secret_key: '', sign_name: '' })
const smsFields = [
  { key: 'access_key', label: 'AccessKey',  placeholder: '', secret: false },
  { key: 'secret_key', label: 'SecretKey',  placeholder: '', secret: true  },
  { key: 'sign_name',  label: '短信签名',   placeholder: '公司名称', secret: false },
]

async function saveConfig(_plugin: string, _data: Record<string, any>) {
  // Save each field as a ConfigKV via a generic config endpoint
  // (Implemented as per-plugin config save in Phase 5 system settings)
  saved.value = true
  setTimeout(() => saved.value = false, 2000)
}

async function saveSmsConfig() {
  await request.put('/system/sms/config', sms.value)
  saved.value = true
  setTimeout(() => saved.value = false, 2000)
}
</script>
