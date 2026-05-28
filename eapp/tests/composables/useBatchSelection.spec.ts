import { describe, expect, it } from 'vitest'
import { useBatchSelection } from '@/composables/useBatchSelection'

describe('useBatchSelection', () => {
  it('starts empty', () => {
    const { count, selected } = useBatchSelection()
    expect(count.value).toBe(0)
    expect(selected.value.size).toBe(0)
  })

  it('toggle is idempotent (add then remove)', () => {
    const { toggle, isSelected, count } = useBatchSelection()
    toggle('a')
    expect(isSelected('a')).toBe(true)
    expect(count.value).toBe(1)
    toggle('a')
    expect(isSelected('a')).toBe(false)
    expect(count.value).toBe(0)
  })

  it('selectAll adds all items and clear empties', () => {
    const { selectAll, clear, count } = useBatchSelection()
    selectAll(['a', 'b', 'c'])
    expect(count.value).toBe(3)
    clear()
    expect(count.value).toBe(0)
  })

  it('isSelected returns correct state', () => {
    const { toggle, isSelected } = useBatchSelection()
    toggle('x')
    toggle('y')
    expect(isSelected('x')).toBe(true)
    expect(isSelected('y')).toBe(true)
    expect(isSelected('z')).toBe(false)
  })

  it('respects max limit', () => {
    const { toggle, count } = useBatchSelection({ max: 2 })
    toggle('a')
    toggle('b')
    toggle('c')
    expect(count.value).toBe(2)
  })
})
