/**
 * 主题归一化：把历史页面/组件里散落的硬编码颜色、圆角统一收敛到
 * App.vue 的设计令牌，避免逐页改样式时遗漏。
 *
 * 用法：node scripts/normalize-theme.cjs
 * 仅处理 pages/ 与 components/ 下的 .vue；刻意保留兼容映射与图标源的文件被排除。
 */
const fs = require('fs')
const path = require('path')

const ROOT = path.join(__dirname, '..')
const DIRS = ['pages', 'components'].map((d) => path.join(ROOT, d))
const EXCLUDE = new Set([
  'MetricCard.vue', // 内含旧色值 → 新变量的兼容映射表
  'StatCard.vue',   // 同上
  'AppIcon.vue',    // 图标组件不含主题色
  'SkeletonCard.vue', // 骨架渐变色属于中性占位色
])

const COLOR_MAP = [
  [/#2563eb/gi, 'var(--eapp-primary)'],
  [/#1d4ed8/gi, 'var(--eapp-primary)'],
  [/#16a34a/gi, 'var(--eapp-success)'],
  [/#f59e0b/gi, 'var(--eapp-warning)'],
  [/#dc2626/gi, 'var(--eapp-danger)'],
  [/#ef4444/gi, 'var(--eapp-danger)'],
  [/#f97316/gi, 'var(--eapp-accent)'],
  [/#eff6ff/gi, 'var(--eapp-primary-soft)'],
  [/#dbeafe/gi, 'var(--eapp-primary-soft)'],
  [/#dcfce7/gi, 'var(--eapp-success-soft)'],
  [/#fef3c7/gi, 'var(--eapp-warning-soft)'],
  [/#fee2e2/gi, 'var(--eapp-danger-soft)'],
  [/#ffedd5/gi, 'var(--eapp-accent-soft)'],
  [/#f8fafc/gi, 'var(--eapp-bg)'],
  [/#f1f5f9/gi, 'var(--eapp-surface-subtle)'],
  [/#e2e8f0/gi, 'var(--eapp-border)'],
  [/#cbd5e1/gi, 'var(--eapp-border-strong)'],
  [/#1e293b/gi, 'var(--eapp-text)'],
  [/#334155/gi, 'var(--eapp-text-secondary)'],
  [/#475569/gi, 'var(--eapp-text-secondary)'],
  [/#64748b/gi, 'var(--eapp-text-muted)'],
  [/#94a3b8/gi, 'var(--eapp-text-faint)'],
  [/#92400e/gi, 'var(--eapp-warning)'],
  [/#166534/gi, 'var(--eapp-success)'],
  [/#991b1b/gi, 'var(--eapp-danger)'],
  [/#c2410c/gi, 'var(--eapp-accent)'],
  // 营销入口历史花色统一收敛到品牌主色
  [/#8b5cf6/gi, 'var(--eapp-primary)'],
  [/#ec4899/gi, 'var(--eapp-primary)'],
  [/#0ea5e9/gi, 'var(--eapp-primary)'],
  [/#d946ef/gi, 'var(--eapp-primary)'],
  [/#14b8a6/gi, 'var(--eapp-primary)'],
]

function walk(dir, out = []) {
  if (!fs.existsSync(dir)) return out
  for (const name of fs.readdirSync(dir)) {
    const p = path.join(dir, name)
    const stat = fs.statSync(p)
    if (stat.isDirectory()) walk(p, out)
    else if (name.endsWith('.vue') && !EXCLUDE.has(name)) out.push(p)
  }
  return out
}

let changedFiles = 0
let totalHits = 0
for (const dir of DIRS) {
  for (const file of walk(dir)) {
    const original = fs.readFileSync(file, 'utf8')
    let next = original
    let hits = 0

    for (const [re, replacement] of COLOR_MAP) {
      const before = next
      next = next.replace(re, replacement)
      if (next !== before) hits++
    }

    // 统一卡片圆角（仅单值写法，不影响 50% / pill / 多值圆角）
    const radiusRe = /border-radius:\s*(?:20|22|24|28|32)rpx\s*;/g
    const withRadius = next.replace(radiusRe, 'border-radius: var(--eapp-radius-lg);')
    if (withRadius !== next) {
      next = withRadius
      hits++
    }

    // 统一描边宽度：小程序端 1px 在高清屏偏粗
    const borderRe = /border:\s*1px solid/g
    const withBorder = next.replace(borderRe, 'border: 1rpx solid')
    if (withBorder !== next) {
      next = withBorder
      hits++
    }

    if (next !== original) {
      fs.writeFileSync(file, next, 'utf8')
      changedFiles++
      totalHits += hits
      console.log(`  ${path.relative(ROOT, file)}`)
    }
  }
}
console.log(`normalized ${changedFiles} files (${totalHits} rule groups applied)`)
