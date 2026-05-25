import { ref } from 'vue'
import type { NotifyPayload } from './notify'

export type ToastItem = Required<Pick<NotifyPayload, 'message' | 'level'>> & {
  id: number
}

const toasts = ref<ToastItem[]>([])
const toastTimers = new Map<number, ReturnType<typeof setTimeout>>()
let toastSeq = 0

export function dismissToast(id: number) {
  const timer = toastTimers.get(id)
  if (timer) {
    clearTimeout(timer)
    toastTimers.delete(id)
  }
  toasts.value = toasts.value.filter((item) => item.id !== id)
}

export function pushToast(payload: NotifyPayload) {
  const message = String(payload.message || '').trim()
  if (!message) return

  toastSeq += 1
  const id = toastSeq
  const level = payload.level || 'info'
  const duration = Number(payload.duration ?? 2200)

  toasts.value.push({ id, message, level })
  if (duration > 0) {
    const timer = setTimeout(() => dismissToast(id), duration)
    toastTimers.set(id, timer)
  }
}

export function useToasts() {
  return {
    toasts,
    pushToast,
    dismissToast,
  }
}
