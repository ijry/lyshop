import { describe, expect, it } from 'vitest'
import { hasPermission } from '@/utils/permission'

describe('hasPermission', () => {
  it('returns true for wildcard', () => {
    expect(hasPermission(['*'], 'order.ship')).toBe(true)
  })

  it('returns false for missing permission', () => {
    expect(hasPermission(['order.list'], 'order.ship')).toBe(false)
  })
})
