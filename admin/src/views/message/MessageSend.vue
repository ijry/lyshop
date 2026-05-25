<template>
  <div>
    <h2 class="text-xl font-semibold text-slate-800 mb-6">{{ $t('message.send.title') }}</h2>
    <div class="bg-white rounded-xl shadow-sm p-6 max-w-xl">
      <div class="space-y-4">
        <div>
          <label class="text-sm font-medium text-slate-700 mb-1 block">{{ $t('message.send.group') }}</label>
          <select v-model="form.group" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500">
            <option value="system">{{ $t('message.group.system') }}</option>
            <option value="order">{{ $t('message.group.order') }}</option>
            <option value="marketing">{{ $t('message.group.marketing') }}</option>
          </select>
        </div>
        <div>
          <label class="text-sm font-medium text-slate-700 mb-1 block">{{ $t('message.send.recipient') }}</label>
          <input v-model.number="form.user_id" type="number" :placeholder="$t('message.send.recipientHint')"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
          <p class="text-xs text-slate-400 mt-1">{{ $t('message.send.recipientNote') }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-slate-700 mb-1 block">{{ $t('message.send.msgTitle') }}</label>
          <input v-model="form.title" :placeholder="$t('message.send.msgTitle')"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
        </div>
        <div>
          <label class="text-sm font-medium text-slate-700 mb-1 block">{{ $t('message.send.msgContent') }}</label>
          <textarea v-model="form.content" rows="4" :placeholder="$t('message.send.contentPlaceholder')"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500 resize-none" />
        </div>
        <button @click="send" :disabled="!form.title || !form.content"
          class="bg-red-600 text-white px-6 py-2.5 rounded-xl text-sm font-medium hover:bg-red-700 transition disabled:opacity-40">
          {{ $t('message.send.submit') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'
import { notify } from '@/utils/notify'

const { t } = useI18n()

const form = ref({ group: 'system', user_id: 0, title: '', content: '' })

async function send() {
  await request.post('/messages/send', form.value)
  notify(t('message.send.success'))
  form.value = { group: 'system', user_id: 0, title: '', content: '' }
}
</script>
