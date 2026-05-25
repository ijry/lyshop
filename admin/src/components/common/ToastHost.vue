<template>
  <div class="fixed top-4 right-4 z-[1000] space-y-2 pointer-events-none">
    <div
      v-for="item in toasts"
      :key="item.id"
      :class="toastClass(item.level)"
      class="min-w-[220px] max-w-[360px] px-3 py-2 rounded-lg shadow-lg border text-sm pointer-events-auto flex items-start gap-2"
    >
      <span class="flex-1 break-words">{{ item.message }}</span>
      <button class="text-xs opacity-70 hover:opacity-100" @click="dismissToast(item.id)">关闭</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { NotifyLevel } from '@/utils/notify'
import { useToasts } from '@/utils/toast'

const { toasts, dismissToast } = useToasts()

function toastClass(level: NotifyLevel) {
  if (level === 'success') return 'bg-emerald-50 border-emerald-200 text-emerald-700'
  if (level === 'warning') return 'bg-amber-50 border-amber-200 text-amber-700'
  if (level === 'error') return 'bg-red-50 border-red-200 text-red-700'
  return 'bg-slate-50 border-slate-200 text-slate-700'
}
</script>
