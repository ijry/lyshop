import { ref } from 'vue'
import {
  getDecorVariants,
  getDecorVariant,
  updateDecorVariant,
  publishDecorVariant,
  copyDecorVariant,
  renameDecorVariant,
  deleteDecorVariant,
} from '@/api/decor'
import { componentLib, createDefaultProps, type DecorComponent } from '@/types/decor'

let compSeq = 100

export function useDecorEditor() {
  const variants = ref<any[]>([])
  const currentVariantKey = ref('default')
  const components = ref<DecorComponent[]>([])
  const selectedIndex = ref(-1)
  const saving = ref(false)

  async function loadVariants() {
    const data: any = await getDecorVariants()
    variants.value = Array.isArray(data) ? data : []
    if (!variants.value.find((v: any) => v.variant_key === currentVariantKey.value) && variants.value.length) {
      currentVariantKey.value = String(variants.value[0].variant_key || 'default')
    }
  }

  async function selectVariant(key: string) {
    currentVariantKey.value = key
    selectedIndex.value = -1
    const data: any = await getDecorVariant(key)
    const raw = data?.components
    try {
      const list = typeof raw === 'string' ? JSON.parse(raw) : raw
      components.value = Array.isArray(list) ? list : []
    } catch {
      components.value = []
    }
  }

  async function save() {
    saving.value = true
    try {
      await updateDecorVariant(currentVariantKey.value, components.value)
      uni.showToast({ title: '保存成功', icon: 'success' })
    } finally {
      saving.value = false
    }
  }

  async function publish() {
    await publishDecorVariant(currentVariantKey.value)
    uni.showToast({ title: '发布成功', icon: 'success' })
    await loadVariants()
  }

  function appendComp(type: string) {
    compSeq += 1
    const comp: DecorComponent = { type, id: `c${compSeq}`, props: createDefaultProps(type) }
    components.value = [...components.value, comp]
    selectedIndex.value = components.value.length - 1
  }

  function moveUp(index: number) {
    if (index <= 0) return
    const arr = [...components.value]
    ;[arr[index - 1], arr[index]] = [arr[index], arr[index - 1]]
    components.value = arr
    if (selectedIndex.value === index) selectedIndex.value = index - 1
    else if (selectedIndex.value === index - 1) selectedIndex.value = index
  }

  function moveDown(index: number) {
    if (index >= components.value.length - 1) return
    const arr = [...components.value]
    ;[arr[index], arr[index + 1]] = [arr[index + 1], arr[index]]
    components.value = arr
    if (selectedIndex.value === index) selectedIndex.value = index + 1
    else if (selectedIndex.value === index + 1) selectedIndex.value = index
  }

  function removeComp(index: number) {
    const arr = [...components.value]
    arr.splice(index, 1)
    components.value = arr
    if (selectedIndex.value === index) selectedIndex.value = -1
    else if (selectedIndex.value > index) selectedIndex.value -= 1
  }

  function selectComp(index: number) {
    selectedIndex.value = index
  }

  return {
    variants, currentVariantKey, components, selectedIndex, saving,
    loadVariants, selectVariant, save, publish,
    appendComp, moveUp, moveDown, removeComp, selectComp,
    componentLib,
    copyDecorVariant, renameDecorVariant, deleteDecorVariant,
  }
}
