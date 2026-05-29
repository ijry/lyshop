import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/decor', () => ({
  getDecorVariants: vi.fn(async () => [
    { variant_key: 'default', variant_name: '默认副本', is_published: true },
  ]),
  getDecorVariant: vi.fn(async () => ({
    variant_key: 'default',
    variant_name: '默认副本',
    components: JSON.stringify([{ type: 'banner', id: 'c1', props: { images: [], height: 350 } }]),
  })),
  updateDecorVariant: vi.fn(async () => null),
  publishDecorVariant: vi.fn(async () => null),
  copyDecorVariant: vi.fn(async () => null),
  renameDecorVariant: vi.fn(async () => null),
  deleteDecorVariant: vi.fn(async () => null),
}))

;(globalThis as any).uni = (globalThis as any).uni || {
  getStorageSync: () => '',
  setStorageSync: () => {},
  removeStorageSync: () => {},
  showToast: () => {},
}

import { useDecorEditor } from '@/composables/useDecorEditor'

describe('useDecorEditor', () => {
  it('loadVariants populates variants list', async () => {
    const editor = useDecorEditor()
    await editor.loadVariants()
    expect(editor.variants.value).toHaveLength(1)
  })
  it('selectVariant parses components', async () => {
    const editor = useDecorEditor()
    await editor.loadVariants()
    await editor.selectVariant('default')
    expect(editor.components.value).toHaveLength(1)
    expect(editor.components.value[0].type).toBe('banner')
  })
  it('appendComp adds a component', async () => {
    const editor = useDecorEditor()
    editor.appendComp('spacer')
    expect(editor.components.value).toHaveLength(1)
    expect(editor.selectedIndex.value).toBe(0)
  })
  it('moveUp and moveDown swap components', () => {
    const editor = useDecorEditor()
    editor.appendComp('banner')
    editor.appendComp('spacer')
    editor.moveDown(0)
    expect(editor.components.value[0].type).toBe('spacer')
    expect(editor.components.value[1].type).toBe('banner')
    editor.moveUp(1)
    expect(editor.components.value[0].type).toBe('banner')
  })
  it('removeComp removes component at index', () => {
    const editor = useDecorEditor()
    editor.appendComp('banner')
    editor.appendComp('spacer')
    editor.removeComp(0)
    expect(editor.components.value).toHaveLength(1)
    expect(editor.components.value[0].type).toBe('spacer')
  })
})
