<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-slate-800">{{ $t('checkin.rules.title') }}</h2>
      <button @click="addRule" class="bg-red-600 text-white px-4 py-2 rounded-xl text-sm hover:bg-red-700 transition">{{ $t('checkin.rules.add') }}</button>
    </div>
    <p class="text-sm text-slate-500 mb-4">{{ $t('checkin.rules.hint') }}</p>

    <div class="bg-white rounded-xl shadow-sm overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 text-slate-500">
          <tr>
            <th class="px-4 py-3 text-left">{{ $t('checkin.rules.day') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('checkin.rules.points') }}</th>
            <th class="px-4 py-3 text-left">{{ $t('common.action') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(r, i) in rules" :key="i" class="border-t border-slate-50">
            <td class="px-4 py-3">
              <input v-model.number="r.day" type="number" min="0" class="w-20 border border-slate-200 rounded-lg px-3 py-1.5 text-sm" />
              <span v-if="r.day === 0" class="text-xs text-slate-400 ml-2">{{ $t('checkin.rules.default') }}</span>
            </td>
            <td class="px-4 py-3">
              <input v-model.number="r.points" type="number" min="1" class="w-24 border border-slate-200 rounded-lg px-3 py-1.5 text-sm" />
            </td>
            <td class="px-4 py-3">
              <button @click="rules.splice(i, 1)" class="text-red-500 hover:text-red-700 text-xs">{{ $t('common.delete') }}</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <button @click="save" class="mt-4 bg-red-600 text-white px-6 py-2.5 rounded-xl text-sm hover:bg-red-700 transition">
      {{ $t('common.save') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'
import { notify } from '@/utils/notify'

const { t } = useI18n()

const rules = ref<any[]>([])

function addRule() {
  rules.value.push({ day: rules.value.length, points: 10 })
}

async function save() {
  await request.put('/checkin/rules', rules.value)
  notify(t('common.saveSuccess'))
}

onMounted(async () => {
  const data: any = await request.get('/checkin/rules')
  rules.value = data || [{ day: 0, points: 10 }]
})
</script>
