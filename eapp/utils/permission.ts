export function parsePermsFromToken(token: string): string[] {
  try {
    const seg = String(token || '').split('.')[1] || ''
    if (!seg) return []
    const normalized = seg.replace(/-/g, '+').replace(/_/g, '/')
    const padded = normalized + '='.repeat((4 - normalized.length % 4) % 4)
    if (typeof atob !== 'function') return []
    const payload = JSON.parse(atob(padded))
    const rows = Array.isArray(payload?.perms) ? payload.perms : []
    return rows.map((row: any) => String(row || ''))
  } catch {
    return []
  }
}

export function hasPermission(perms: string[], permission: string): boolean {
  if (!permission) return true
  if (!Array.isArray(perms) || perms.length === 0) return false
  return perms.includes('*') || perms.includes(permission)
}
