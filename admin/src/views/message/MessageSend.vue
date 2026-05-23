<template>
  <div>
    <h2 class="text-xl font-semibold text-slate-800 mb-6">发送消息</h2>
    <div class="bg-white rounded-xl shadow-sm p-6 max-w-xl">
      <div class="space-y-4">
        <div>
          <label class="text-sm font-medium text-slate-700 mb-1 block">消息分组</label>
          <select v-model="form.group" class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500">
            <option value="system">系统通知</option>
            <option value="order">订单消息</option>
            <option value="marketing">营销消息</option>
          </select>
        </div>
        <div>
          <label class="text-sm font-medium text-slate-700 mb-1 block">接收用户</label>
          <input v-model.number="form.user_id" type="number" placeholder="用户ID（0=全部用户）"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
          <p class="text-xs text-slate-400 mt-1">填 0 表示广播给所有用户</p>
        </div>
        <div>
          <label class="text-sm font-medium text-slate-700 mb-1 block">消息标题</label>
          <input v-model="form.title" placeholder="消息标题"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500" />
        </div>
        <div>
          <label class="text-sm font-medium text-slate-700 mb-1 block">消息内容</label>
          <textarea v-model="form.content" rows="4" placeholder="消息正文内容"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-red-500 resize-none" />
        </div>
        <button @click="send" :disabled="!form.title || !form.content"
          class="bg-red-600 text-white px-6 py-2.5 rounded-xl text-sm font-medium hover:bg-red-700 transition disabled:opacity-40">
          发送消息
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import request from '@/api/request'

const form = ref({ group: 'system', user_id: 0, title: '', content: '' })

async function send() {
  await request.post('/messages/send', form.value)
  alert('发送成功')
  form.value = { group: 'system', user_id: 0, title: '', content: '' }
}
</script>
