export type NotifyLevel = 'info' | 'success' | 'warning' | 'error'

export type NotifyPayload = {
  message: string
  level?: NotifyLevel
  duration?: number
}

export type NotifyHandler = (payload: NotifyPayload) => void

function defaultNotifyHandler(payload: NotifyPayload) {
  if (!payload.message) return
  if (typeof window !== 'undefined' && typeof window.alert === 'function') {
    window.alert(payload.message)
  }
}

let activeNotifyHandler: NotifyHandler = defaultNotifyHandler

export function setNotifyHandler(handler?: NotifyHandler | null) {
  activeNotifyHandler = handler || defaultNotifyHandler
}

export function notify(message: string, options?: Omit<NotifyPayload, 'message'>) {
  if (!message) return
  activeNotifyHandler({
    message,
    level: options?.level || 'info',
    duration: options?.duration,
  })
}
