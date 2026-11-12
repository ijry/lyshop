/**
 * 商家移动端统一线性图标库
 *
 * 设计约定：24x24 viewBox、1.7 描边、round 端点，无填充、无渐变、无 emoji。
 * 运行时拼装为 base64 SVG data URI，供 <image> / background-image 使用，
 * 同时兼容 H5、微信小程序与 App（避免内联 <svg> 标签在小程序端不可用的问题）。
 */

export const ICON_PATHS = {
  /* 基础 / 导航 */
  circle: '<circle cx="12" cy="12" r="9"/>',
  'chevron-right': '<path d="M9 5l7 7-7 7"/>',
  'chevron-left': '<path d="M15 5l-7 7 7 7"/>',
  'chevron-down': '<path d="M5 9l7 7 7-7"/>',
  'arrow-right': '<path d="M4 12h16M14 6l6 6-6 6"/>',
  'arrow-left': '<path d="M20 12H4M10 6l-6 6 6 6"/>',
  'arrow-up': '<path d="M12 20V4M6 10l6-6 6 6"/>',
  'arrow-down': '<path d="M12 4v16M6 14l6 6 6-6"/>',
  check: '<path d="M5 12.5l4.5 4.5L19 7"/>',
  close: '<path d="M6 6l12 12M18 6L6 18"/>',
  plus: '<path d="M12 5v14M5 12h14"/>',
  minus: '<path d="M5 12h14"/>',
  more: '<path d="M5 12h.01M12 12h.01M19 12h.01" stroke-width="2.6"/>',
  refresh: '<path d="M20 11a8 8 0 1 0-2.3 6.3M20 5v6h-6"/>',

  /* 业务模块 */
  dashboard:
    '<rect x="3" y="3" width="7.5" height="7.5" rx="1.5"/><rect x="13.5" y="3" width="7.5" height="7.5" rx="1.5"/><rect x="3" y="13.5" width="7.5" height="7.5" rx="1.5"/><rect x="13.5" y="13.5" width="7.5" height="7.5" rx="1.5"/>',
  order:
    '<path d="M5 3h14v18l-2.5-1.5L14 21l-2-1.5L10 21l-2.5-1.5L5 21V3z"/><path d="M9 8h6M9 12h6M9 16h4"/>',
  package:
    '<path d="M3 7.5L12 3l9 4.5v9L12 21l-9-4.5v-9z"/><path d="M3 7.5L12 12l9-4.5M12 12v9"/>',
  product:
    '<path d="M20.5 13.4l-7.1 7.1a2 2 0 0 1-2.8 0l-7-7a2 2 0 0 1 0-2.8l7.1-7.1a2 2 0 0 1 1.4-.6H19a1.5 1.5 0 0 1 1.5 1.5v7.5a2 2 0 0 1 0 1.4z"/><circle cx="16.5" cy="7.5" r="1.2"/>',
  'after-sale':
    '<path d="M20 7a5 5 0 0 1-6.4 6.3L7 20a2 2 0 0 1-3-3l6.7-6.5A5 5 0 0 1 17 4l-3.2 3.2.5 2.6 2.6.5L20 7z"/>',
  star: '<path d="M12 3l2.7 5.6 6.1.9-4.4 4.3 1 6.1-5.4-2.9-5.4 2.9 1-6.1L3.2 9.5l6.1-.9L12 3z"/>',
  decor:
    '<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18M9 9v12M15 13h3"/>',
  chart:
    '<path d="M3 20V4M3 20h18"/><rect x="7" y="12" width="3" height="5"/><rect x="12" y="8" width="3" height="9"/><rect x="17" y="5" width="3" height="12"/>',
  'chart-line': '<path d="M3 20h18M4 16l5-5 4 3 6-8"/><path d="M17 6h2v2"/>',
  warehouse:
    '<path d="M3 9l9-6 9 6v11H3V9z"/><path d="M3 9h18M8 20v-6h8v6"/>',
  coupon:
    '<path d="M3 6h18v4a2 2 0 0 0 0 4v4H3v-4a2 2 0 0 0 0-4V6z"/><path d="M12 7v2M12 12v2M12 17v1"/>',
  flash: '<path d="M13 2L5 13h6l-1 9 8-11h-6l1-9z"/>',
  users:
    '<circle cx="9" cy="8" r="3"/><path d="M3 20v-1a6 6 0 0 1 12 0v1"/><path d="M16.5 5.5a3 3 0 0 1 0 5.5M18 14a5 5 0 0 1 3 5v1"/>',
  money:
    '<rect x="2" y="6" width="20" height="12" rx="2"/><circle cx="12" cy="12" r="2.5"/><path d="M6 12h.01M18 12h.01"/>',
  crown: '<path d="M3 8l4 4 5-7 5 7 4-4-2 11H5L3 8z"/><path d="M5 19h14"/>',
  megaphone:
    '<path d="M4 10v4h3l8 4V6l-8 4H4z"/><path d="M18 9a4 4 0 0 1 0 6M7 14l1 6h3l-1-6"/>',
  clipboard:
    '<rect x="5" y="4" width="14" height="17" rx="2"/><path d="M9 4V3h6v1M9 10h6M9 14h6M9 18h4"/>',
  share:
    '<circle cx="6" cy="12" r="2.5"/><circle cx="18" cy="6" r="2.5"/><circle cx="18" cy="18" r="2.5"/><path d="M8.3 10.8l7.4-3.6M8.3 13.2l7.4 3.6"/>',
  diamond:
    '<path d="M6 4h12l4 6-10 11L2 10l4-6z"/><path d="M2 10h20M6 4l4 6 2 11 2-11 4-6"/>',

  /* 操作 / 状态 */
  scan:
    '<path d="M4 8V5a1 1 0 0 1 1-1h3M16 4h3a1 1 0 0 1 1 1v3M20 16v3a1 1 0 0 1-1 1h-3M8 20H5a1 1 0 0 1-1-1v-3M3 12h18"/>',
  chat:
    '<path d="M20 15a3 3 0 0 1-3 3H9l-5 3v-3a3 3 0 0 1-1-2V7a3 3 0 0 1 3-3h11a3 3 0 0 1 3 3v8z"/><path d="M8 10h8M8 13h5"/>',
  headset:
    '<path d="M4 14v-2a8 8 0 0 1 16 0v2"/><rect x="2" y="13" width="4" height="7" rx="1.5"/><rect x="18" y="13" width="4" height="7" rx="1.5"/><path d="M20 20v1a3 3 0 0 1-3 3h-3"/>',
  bell:
    '<path d="M18 9a6 6 0 1 0-12 0c0 5-2 6-2 7h16c0-1-2-2-2-7z"/><path d="M10 20a2 2 0 0 0 4 0"/>',
  filter: '<path d="M3 5h18M6 12h12M10 19h4"/>',
  search: '<circle cx="10.5" cy="10.5" r="6.5"/><path d="M15.5 15.5L21 21"/>',
  user: '<circle cx="12" cy="8" r="3.5"/><path d="M5 20a7 7 0 0 1 14 0"/>',
  settings:
    '<circle cx="12" cy="12" r="3"/><path d="M12 2.5v3M12 18.5v3M2.5 12h3M18.5 12h3M5.2 5.2l2.1 2.1M16.7 16.7l2.1 2.1M18.8 5.2l-2.1 2.1M7.3 16.7l-2.1 2.1"/>',
  shield:
    '<path d="M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6l8-3z"/><path d="M9 12l2 2 4-4"/>',
  printer:
    '<path d="M6 9V3h12v6M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect x="6" y="15" width="12" height="6" rx="1"/>',
  truck:
    '<path d="M2 6h12v11H2zM14 10h4l4 4v3h-8"/><circle cx="6" cy="18" r="1.8"/><circle cx="18" cy="18" r="1.8"/>',
  edit: '<path d="M4 20h4l11-11-4-4L4 16v4z"/><path d="M14 6l4 4"/>',
  clock: '<circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/>',
  calendar:
    '<rect x="3" y="5" width="18" height="16" rx="2"/><path d="M3 10h18M8 3v4M16 3v4"/>',
  home: '<path d="M3 10l9-7 9 7v10a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V10z"/><path d="M9 21v-7h6v7"/>',
  inbox:
    '<path d="M3 13l3-8h12l3 8v6a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1v-6z"/><path d="M3 13h5l1 3h6l1-3h5"/>',
  document:
    '<path d="M6 3h8l5 5v12a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z"/><path d="M14 3v5h5M9 13h7M9 17h5"/>',
  wallet: '<rect x="3" y="6" width="18" height="14" rx="2"/><path d="M3 10h18M16 15h2"/>',
  store:
    '<path d="M4 9V5h16v4M3 9h18l-1 3H4L3 9zM5 12v8h14v-8M10 20v-5h4v5"/>',
  logout: '<path d="M10 4H5a1 1 0 0 0-1 1v14a1 1 0 0 0 1 1h5"/><path d="M15 8l4 4-4 4M19 12H9"/>',
  info: '<circle cx="12" cy="12" r="9"/><path d="M12 11v5M12 8h.01"/>',
  warning:
    '<path d="M12 3l9.5 17H2.5L12 3z"/><path d="M12 9v5M12 17h.01"/>',
  'check-circle': '<circle cx="12" cy="12" r="9"/><path d="M8 12l3 3 5-6"/>',
  'close-circle': '<circle cx="12" cy="12" r="9"/><path d="M9 9l6 6M15 9l-6 6"/>',
  eye: '<path d="M2 12s3.5-6 10-6 10 6 10 6-3.5 6-10 6-10-6-10-6z"/><circle cx="12" cy="12" r="2.5"/>',
  lock: '<rect x="4" y="10" width="16" height="11" rx="2"/><path d="M8 10V7a4 4 0 0 1 8 0v3"/>',
  mail: '<rect x="3" y="5" width="18" height="14" rx="2"/><path d="M3 7l9 6 9-6"/>',
  list: '<path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01" stroke-width="2"/>',
  layers: '<path d="M12 3l9 5-9 5-9-5 9-5z"/><path d="M3 13l9 5 9-5M3 17l9 5 9-5"/>',
  pie: '<path d="M12 3a9 9 0 1 0 9 9h-9V3z"/><path d="M15 3.5A9 9 0 0 1 20.5 9H15V3.5z"/>',
  'trend-up': '<path d="M3 17l6-6 4 4 8-8M15 7h6v6"/>',
  'trend-down': '<path d="M3 7l6 6 4-4 8 8M15 17h6v-6"/>',
  help: '<circle cx="12" cy="12" r="9"/><path d="M9.5 9a2.5 2.5 0 1 1 3.5 2.3c-.8.4-1 1-1 1.7M12 17h.01"/>',
  tag: '<path d="M3 11V4a1 1 0 0 1 1-1h7l10 10-8 8L3 11z"/><circle cx="7.5" cy="7.5" r="1.2"/>',
  target: '<circle cx="12" cy="12" r="9"/><circle cx="12" cy="12" r="5"/><circle cx="12" cy="12" r="1.2"/>',
} as const

export type IconName = keyof typeof ICON_PATHS

/** 语义色键 → 具体色值（SVG data URI 内无法解析 CSS 变量，故使用固定色值） */
export const ICON_COLORS = {
  primary: '#c62f24',
  'primary-deep': '#8f2018',
  text: '#101828',
  secondary: '#475467',
  muted: '#667085',
  faint: '#98a2b3',
  success: '#157f4b',
  warning: '#a75a08',
  danger: '#b42318',
  accent: '#0e7490',
  white: '#ffffff',
} as const

/** 兼容历史代码里直接传 CSS 变量的写法 */
const CSS_VAR_MAP: Record<string, string> = {
  'var(--eapp-primary)': ICON_COLORS.primary,
  'var(--eapp-text)': ICON_COLORS.text,
  'var(--eapp-text-secondary)': ICON_COLORS.secondary,
  'var(--eapp-text-muted)': ICON_COLORS.muted,
  'var(--eapp-text-faint)': ICON_COLORS.faint,
  'var(--eapp-success)': ICON_COLORS.success,
  'var(--eapp-warning)': ICON_COLORS.warning,
  'var(--eapp-danger)': ICON_COLORS.danger,
  'var(--eapp-accent)': ICON_COLORS.accent,
}

export function resolveIconColor(input?: string): string {
  if (!input) return ICON_COLORS.secondary
  if (CSS_VAR_MAP[input]) return CSS_VAR_MAP[input]
  if (input in ICON_COLORS) return ICON_COLORS[input as keyof typeof ICON_COLORS]
  if (input.startsWith('#') || input.startsWith('rgb')) return input
  return ICON_COLORS.secondary
}

const B64_CHARS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'

/** 纯 ASCII base64 编码（SVG 内容全部为 ASCII，小程序端无 btoa） */
function toBase64(input: string): string {
  let output = ''
  for (let i = 0; i < input.length; i += 3) {
    const c1 = input.charCodeAt(i)
    const has2 = i + 1 < input.length
    const has3 = i + 2 < input.length
    const c2 = has2 ? input.charCodeAt(i + 1) : 0
    const c3 = has3 ? input.charCodeAt(i + 2) : 0
    const e1 = c1 >> 2
    const e2 = ((c1 & 3) << 4) | (c2 >> 4)
    const e3 = has2 ? ((c2 & 15) << 2) | (c3 >> 6) : 64
    const e4 = has3 ? c3 & 63 : 64
    output +=
      B64_CHARS[e1] +
      B64_CHARS[e2] +
      (e3 === 64 ? '=' : B64_CHARS[e3]) +
      (e4 === 64 ? '=' : B64_CHARS[e4])
  }
  return output
}

/**
 * 生成图标 data URI。
 * @param name       图标名，未知时降级为 circle
 * @param color      语义色键 / #hex / CSS 变量字符串
 * @param strokeWidth 描边宽度，默认 1.7
 */
export function iconDataUri(name: string, color?: string, strokeWidth = 1.7): string {
  const body = ICON_PATHS[name as IconName] || ICON_PATHS.circle
  const stroke = resolveIconColor(color)
  const svg =
    `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" ` +
    `fill="none" stroke="${stroke}" stroke-width="${strokeWidth}" ` +
    `stroke-linecap="round" stroke-linejoin="round">${body}</svg>`
  return 'data:image/svg+xml;base64,' + toBase64(svg)
}
