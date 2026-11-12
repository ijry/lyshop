import { describe, expect, it } from 'vitest'
import {
  ACTIVITY_FILTERS,
  ACTIVITY_KINDS,
  ACTIVITY_META,
  MARKETING_ROUTES,
  activityDetailUrl,
  activityMeta,
  activityProgress,
  activityState,
  formatActivityTime,
  formatMoney,
  parseSkuLabel,
} from '@/utils/marketing'

describe('marketing metadata', () => {
  it('defines metadata for every activity kind', () => {
    expect(ACTIVITY_KINDS).toEqual(['seckill', 'group-buy', 'bargain'])
    for (const kind of ACTIVITY_KINDS) {
      const meta = ACTIVITY_META[kind]
      expect(meta.kind).toBe(kind)
      expect(meta.title).toBeTruthy()
      expect(meta.iconName).toBeTruthy()
      expect(meta.route).toBe(MARKETING_ROUTES[kind === 'group-buy' ? 'groupBuy' : kind])
    }
  })

  it('falls back to seckill metadata for unknown kinds', () => {
    expect(activityMeta('unknown').kind).toBe('seckill')
    expect(activityMeta(null).kind).toBe('seckill')
  })

  it('builds an encoded detail url', () => {
    const url = activityDetailUrl('seckill', 12, '午场 秒杀')
    expect(url).toContain('kind=seckill')
    expect(url).toContain('id=12')
    expect(url).toContain(encodeURIComponent('午场 秒杀'))
  })

  it('offers a filter for every phase', () => {
    expect(ACTIVITY_FILTERS.map((item) => item.key)).toEqual([
      'all',
      'running',
      'upcoming',
      'ended',
      'disabled',
    ])
  })
})

describe('activity state model', () => {
  const now = new Date('2026-06-01T12:00:00Z').getTime()
  const base = {
    start_at: '2026-06-01T10:00:00Z',
    end_at: '2026-06-01T14:00:00Z',
    status: 1,
  }

  it('derives running / upcoming / ended from the clock', () => {
    expect(activityState(base, now).phase).toBe('running')
    expect(activityState(base, now - 3 * 3600_000).phase).toBe('upcoming')
    expect(activityState(base, now + 3 * 3600_000).phase).toBe('ended')
  })

  it('always reports disabled when status is off', () => {
    expect(activityState({ ...base, status: 0 }, now).label).toBe('已停用')
  })

  it('reports pending schedule when the period is missing', () => {
    expect(activityState({ status: 1 }, now).label).toBe('待排期')
  })

  it('computes linear progress and clamps to 0-100', () => {
    expect(activityProgress(base, now)).toBe(50)
    expect(activityProgress(base, now - 10 * 3600_000)).toBe(0)
    expect(activityProgress(base, now + 10 * 3600_000)).toBe(100)
  })
})

describe('marketing display helpers', () => {
  it('formats ISO time in local time', () => {
    const d = new Date(2026, 4, 27, 10, 30)
    expect(formatActivityTime(d.toISOString(), true)).toBe('2026-05-27 10:30')
    expect(formatActivityTime('')).toBe('-')
  })

  it('formats money with two decimals', () => {
    expect(formatMoney(3599)).toBe('¥3599.00')
    expect(formatMoney(undefined)).toBe('¥0.00')
  })

  it('parses SKU attrs from array, JSON string and falls back to id', () => {
    const attrs = [{ name: '颜色', value: '曜石黑' }, { name: '尺码', value: 'M' }]
    expect(parseSkuLabel(attrs, 7)).toBe('颜色:曜石黑 / 尺码:M')
    expect(parseSkuLabel(JSON.stringify(attrs), 7)).toBe('颜色:曜石黑 / 尺码:M')
    expect(parseSkuLabel('not-json', 7)).toBe('SKU #7')
  })
})
