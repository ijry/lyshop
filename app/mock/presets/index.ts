import type { MockPreset } from './types'
import { mall } from './mall'
import { supermarket } from './supermarket'
import { fresh } from './fresh'
import { jewelry } from './jewelry'
import { farm } from './farm'
import { cake } from './cake'
import { mother } from './mother'

const presets: Record<string, MockPreset> = {
  mall, supermarket, fresh, jewelry, farm, cake, mother,
}

function resolvePresetKey(): string {
  if (typeof window !== 'undefined' && window.location?.search) {
    const params = new URLSearchParams(window.location.search)
    const key = params.get('demo')
    if (key && key in presets) return key
  }
  return 'mall'
}

let resolved: MockPreset | null = null

export function getPreset(): MockPreset {
  if (!resolved) {
    resolved = presets[resolvePresetKey()]
  }
  return resolved
}

export function getPresetKey(): string {
  return getPreset().key
}

export function listPresets(): Array<{ key: string; name: string }> {
  return Object.values(presets).map(p => ({ key: p.key, name: p.name }))
}

export type { MockPreset }
