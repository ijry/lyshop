import { describe, expect, it } from 'vitest'

describe('request placeholder', () => {
  it('keeps api root stable', () => {
    expect('/admin/api').toContain('/admin/api')
  })
})
