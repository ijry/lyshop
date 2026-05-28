import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useBadgeStore } from '@/stores/badge'

describe('badge store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('syncs from dashboard data', () => {
    const store = useBadgeStore()
    store.syncFromDashboard({ pending_ship: 2, pending_after_sale: 1, unread_message: 4 })
    expect(store.orderBadge).toBe(3)
    expect(store.messageBadge).toBe(4)
  })
})
