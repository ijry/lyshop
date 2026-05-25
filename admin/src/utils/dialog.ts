export function confirmAction(message: string): boolean {
  if (typeof window === 'undefined' || typeof window.confirm !== 'function') {
    return false
  }
  return window.confirm(message)
}

export function promptText(message: string, defaultValue = ''): string | null {
  if (typeof window === 'undefined' || typeof window.prompt !== 'function') {
    return null
  }
  return window.prompt(message, defaultValue)
}
