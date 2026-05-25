<template>
  <div class="max-w-5xl mx-auto px-6 py-8">
    <div class="flex flex-col lg:flex-row gap-6">
      <!-- Sidebar -->
      <div class="w-72 shrink-0">
        <!-- Profile card -->
        <div class="card p-6 mb-4">
          <div class="flex items-center gap-4 mb-4">
            <img :src="user.avatar || 'https://api.dicebear.com/7.x/adventurer/svg?seed=default'"
              class="w-16 h-16 rounded-full border-2 border-red-100" />
            <div>
              <h2 class="text-lg font-bold text-gray-900">{{ user.nickname || $t('userCenter.notLoggedIn') }}</h2>
              <p class="text-sm text-gray-500">{{ user.phone || '' }}</p>
            </div>
          </div>
          <div class="bg-gradient-to-r from-red-50 to-orange-50 rounded-xl p-4">
            <div class="flex-between">
              <span class="text-sm text-gray-600">{{ $t('userCenter.myPoints') }}</span>
              <span class="text-2xl font-bold text-red-500">{{ user.points || 0 }}</span>
            </div>
            <p class="text-xs text-gray-400 mt-1">{{ $t('userCenter.pointsExchange') }}</p>
          </div>
        </div>

        <!-- Menu -->
        <div class="card overflow-hidden">
          <button v-for="item in menuItems" :key="item.key"
            @click="item.action"
            :class="activeMenu === item.key ? 'bg-red-50 text-red-600 font-medium border-l-3 border-red-600' : 'text-gray-600 hover:bg-gray-50'"
            class="w-full text-left px-5 py-3.5 text-sm transition-colors flex-between">
            <div class="flex items-center gap-3">
              <div :class="item.icon" class="text-lg" />
              <span>{{ item.label }}</span>
            </div>
            <div class="i-carbon-chevron-right text-xs text-gray-300" />
          </button>
        </div>
      </div>

      <!-- Main content -->
      <div class="flex-1 min-w-0">
        <!-- Order status shortcuts -->
        <div v-if="activeMenu === 'overview'" class="card p-6 mb-4">
          <div class="flex-between mb-4">
            <h3 class="text-base font-bold text-gray-900">{{ $t('userCenter.myOrders') }}</h3>
            <router-link to="/orders" class="text-sm text-red-600 hover:text-red-700">{{ $t('userCenter.allOrders') }}</router-link>
          </div>
          <div class="grid grid-cols-4 gap-4">
            <router-link v-for="s in orderStatuses" :key="s.label" to="/orders"
              class="flex flex-col items-center gap-2 p-4 rounded-xl hover:bg-gray-50 transition-colors">
              <div :class="s.icon" class="text-2xl text-gray-600" />
              <span class="text-sm text-gray-600">{{ s.label }}</span>
            </router-link>
          </div>
        </div>

        <!-- Coupons -->
        <div v-if="activeMenu === 'overview' || activeMenu === 'coupons'" class="card p-6 mb-4">
          <h3 class="text-base font-bold text-gray-900 mb-4">{{ $t('userCenter.myCoupons') }}</h3>
          <div v-if="coupons.length" class="space-y-3">
            <div v-for="c in coupons" :key="c.id"
              class="flex items-center border rounded-xl overflow-hidden"
              :class="c.status === 1 ? 'border-red-200' : 'border-gray-100 opacity-60'">
              <div class="w-24 py-4 text-center shrink-0"
                :class="c.status === 1 ? 'bg-red-500 text-white' : 'bg-gray-200 text-gray-500'">
                <div class="text-xl font-bold">¥{{ c.discount || (c.type === 2 ? (c.discount_rate * 10) + t('userCenter.couponDiscount') : '?') }}</div>
                <div class="text-xs mt-0.5">{{ c.min_amount > 0 ? `${$t('userCenter.couponThreshold')}${c.min_amount}${$t('userCenter.couponAvailable')}` : $t('userCenter.couponNoThreshold') }}</div>
              </div>
              <div class="flex-1 px-4 py-3">
                <p class="text-sm font-medium text-gray-800">{{ c.name || $t('userCenter.myCoupons') }}</p>
                <p class="text-xs text-gray-400 mt-1">{{ c.status === 1 ? $t('userCenter.unused') : c.status === 2 ? $t('userCenter.used') : $t('userCenter.expired') }}</p>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-400 text-sm">{{ $t('userCenter.noCoupons') }}</div>
        </div>

        <!-- Favorites -->
        <div v-if="activeMenu === 'overview' || activeMenu === 'favorites'" class="card p-6 mb-4">
          <div class="flex-between mb-4">
            <h3 class="text-base font-bold text-gray-900">{{ $t('userCenter.myFavorites') }}</h3>
            <span class="text-xs text-gray-400">{{ $t('userCenter.total') }} {{ favoriteTotal }} {{ $t('userCenter.unit') }}</span>
          </div>
          <div v-if="favorites.length" class="space-y-3">
            <div v-for="item in favorites" :key="item.id" class="border border-gray-100 rounded-xl p-3 flex gap-3">
              <img :src="item.cover" class="w-16 h-16 rounded-lg object-cover border border-gray-100 cursor-pointer" @click="toProductDetail(item.id)" />
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-800 font-medium truncate cursor-pointer" @click="toProductDetail(item.id)">{{ item.title }}</p>
                <p class="text-xs text-gray-400 mt-1">{{ $t('productDetail.favorite') }} {{ item.favorite_count || 0 }}</p>
                <div class="flex-between mt-2">
                  <span class="text-sm text-red-500 font-semibold">¥{{ item.price }}</span>
                  <button class="text-xs text-red-500 hover:underline" @click="unfavorite(item.id)">{{ $t('userCenter.unfavorite') }}</button>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-400 text-sm">{{ $t('userCenter.noFavorites') }}</div>
        </div>

        <!-- Points log -->
        <div v-if="activeMenu === 'overview' || activeMenu === 'points'" class="card p-6 mb-4">
          <h3 class="text-base font-bold text-gray-900 mb-4">{{ $t('userCenter.pointsDetail') }}</h3>
          <div v-if="pointsLogs.length" class="divide-y divide-gray-50">
            <div v-for="log in pointsLogs" :key="log.id" class="flex-between py-3">
              <div>
                <p class="text-sm text-gray-800">{{ log.remark }}</p>
                <p class="text-xs text-gray-400 mt-0.5">{{ log.created_at?.slice(0, 10) }}</p>
              </div>
              <span class="text-sm font-bold" :class="log.points > 0 ? 'text-green-600' : 'text-red-500'">
                {{ log.points > 0 ? '+' : '' }}{{ log.points }}
              </span>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-400 text-sm">{{ $t('userCenter.noPointsRecord') }}</div>
        </div>

        <!-- Addresses -->
        <div v-if="activeMenu === 'overview' || activeMenu === 'addresses'" class="card p-6">
          <div class="flex-between mb-4">
            <h3 class="text-base font-bold text-gray-900">{{ $t('userCenter.addressTitle') }}</h3>
            <button class="btn-primary text-xs !px-4 !py-1.5" @click="openAddressEditor()">{{ $t('userCenter.addAddress') }}</button>
          </div>
          <div v-if="addresses.length" class="space-y-3">
            <div v-for="addr in addresses" :key="addr.id"
              class="border border-gray-100 rounded-xl p-4 hover:border-gray-200 transition-colors">
              <div class="flex-between mb-2">
                <div class="flex items-center gap-3">
                  <span class="text-sm font-semibold text-gray-800">{{ addr.name }}</span>
                  <span class="text-sm text-gray-500">{{ addr.phone }}</span>
                  <span v-if="addr.is_default"
                    class="text-xs bg-red-50 text-red-600 px-2 py-0.5 rounded">{{ $t('userCenter.defaultTag') }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <button class="text-xs text-blue-600 hover:underline" @click="openAddressEditor(addr)">{{ $t('userCenter.edit') }}</button>
                  <button class="text-xs text-red-500 hover:underline" @click="removeAddress(addr)">{{ $t('userCenter.delete') }}</button>
                </div>
              </div>
              <p class="text-sm text-gray-500">{{ addr.province }}{{ addr.city }}{{ addr.district }} {{ addr.detail }}</p>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-400 text-sm">{{ $t('userCenter.noAddress') }}</div>
        </div>
      </div>
    </div>

    <div v-if="showAddressDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/35 p-4" @click.self="closeAddressEditor">
      <div class="w-full max-w-lg rounded-2xl bg-white shadow-xl p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ editingAddressID ? $t('userCenter.editAddress') : $t('userCenter.addAddress') }}</h3>
        <div class="space-y-3">
          <input v-model.trim="addressForm.name" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-red-200" :placeholder="$t('userCenter.recipientName')" />
          <input v-model.trim="addressForm.phone" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-red-200" :placeholder="$t('userCenter.phone')" maxlength="11" />
          <div class="grid grid-cols-3 gap-3">
            <input v-model.trim="addressForm.province" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-red-200" :placeholder="$t('userCenter.province')" />
            <input v-model.trim="addressForm.city" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-red-200" :placeholder="$t('userCenter.city')" />
            <input v-model.trim="addressForm.district" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-red-200" :placeholder="$t('userCenter.district')" />
          </div>
          <input v-model.trim="addressForm.detail" class="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-red-200" :placeholder="$t('userCenter.detailAddress')" />
          <label class="flex items-center gap-2 text-sm text-gray-600">
            <input type="checkbox" v-model="addressForm.is_default" />
            {{ $t('userCenter.setDefault') }}
          </label>
        </div>
        <div class="flex justify-end gap-2 mt-5">
          <button class="btn-outline !px-5 !py-2 text-xs" @click="closeAddressEditor">{{ $t('userCenter.cancel') }}</button>
          <button class="btn-primary !px-5 !py-2 text-xs" :disabled="savingAddress" @click="saveAddress">{{ $t('userCenter.save') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { del, get, post, put } from '@/api/request'

const { t } = useI18n()
const user = ref<any>({})
const coupons = ref<any[]>([])
const pointsLogs = ref<any[]>([])
const addresses = ref<any[]>([])
const favorites = ref<any[]>([])
const favoriteTotal = ref(0)
const activeMenu = ref('overview')
const router = useRouter()
const showAddressDialog = ref(false)
const savingAddress = ref(false)
const editingAddressID = ref(0)
const addressForm = ref({
  name: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  is_default: false,
})

const menuItems = computed(() => [
  { key: 'overview', label: t('userCenter.myOrders'), icon: 'i-carbon-dashboard', action: () => activeMenu.value = 'overview' },
  { key: 'coupons', label: t('userCenter.myCoupons'), icon: 'i-carbon-ticket', action: () => activeMenu.value = 'coupons' },
  { key: 'favorites', label: t('userCenter.myFavorites'), icon: 'i-carbon-favorite', action: () => activeMenu.value = 'favorites' },
  { key: 'points', label: t('userCenter.myPoints'), icon: 'i-carbon-star', action: () => activeMenu.value = 'points' },
  { key: 'addresses', label: t('userCenter.addressTitle'), icon: 'i-carbon-location', action: () => activeMenu.value = 'addresses' },
])

const orderStatuses = computed(() => [
  { label: t('userCenter.unpaid'), icon: 'i-carbon-wallet' },
  { label: t('userCenter.unshipped'), icon: 'i-carbon-package' },
  { label: t('userCenter.unreceived'), icon: 'i-carbon-delivery-truck' },
  { label: t('userCenter.afterSale'), icon: 'i-carbon-renew' },
])

onMounted(async () => {
  const [profile, couponData, pointsData, addrData, favoriteData] = await Promise.all([
    get<any>('/api/v1/user/profile'),
    get<any>('/api/v1/user/coupons'),
    get<any>('/api/v1/user/points/logs'),
    get<any[]>('/api/v1/addresses'),
    get<any>('/api/v1/user/favorites', { page: 1, size: 12 }),
  ])
  user.value = profile || {}
  coupons.value = couponData || []
  pointsLogs.value = pointsData?.list || []
  addresses.value = addrData || []
  favorites.value = Array.isArray(favoriteData?.list) ? favoriteData.list : []
  favoriteTotal.value = Number(favoriteData?.total || 0)
})

function resetAddressForm() {
  editingAddressID.value = 0
  addressForm.value = {
    name: '',
    phone: '',
    province: '',
    city: '',
    district: '',
    detail: '',
    is_default: false,
  }
}

function openAddressEditor(addr?: any) {
  if (!addr) {
    resetAddressForm()
    showAddressDialog.value = true
    return
  }
  editingAddressID.value = Number(addr.id || 0)
  addressForm.value = {
    name: String(addr.name || ''),
    phone: String(addr.phone || ''),
    province: String(addr.province || ''),
    city: String(addr.city || ''),
    district: String(addr.district || ''),
    detail: String(addr.detail || ''),
    is_default: Number(addr.is_default || 0) === 1,
  }
  showAddressDialog.value = true
}

function closeAddressEditor() {
  showAddressDialog.value = false
  resetAddressForm()
}

function validateAddress() {
  if (!addressForm.value.name.trim()) return t('userCenter.validateName')
  if (!/^1\d{10}$/.test(addressForm.value.phone)) return t('userCenter.validatePhone')
  if (!addressForm.value.province.trim() || !addressForm.value.city.trim() || !addressForm.value.district.trim()) return t('userCenter.validateRegion')
  if (!addressForm.value.detail.trim()) return t('userCenter.validateDetail')
  return ''
}

async function refreshAddresses() {
  addresses.value = await get<any[]>('/api/v1/addresses') || []
}

async function saveAddress() {
  const error = validateAddress()
  if (error) {
    alert(error)
    return
  }
  if (savingAddress.value) return
  savingAddress.value = true
  const payload = {
    name: addressForm.value.name.trim(),
    phone: addressForm.value.phone.trim(),
    province: addressForm.value.province.trim(),
    city: addressForm.value.city.trim(),
    district: addressForm.value.district.trim(),
    detail: addressForm.value.detail.trim(),
    is_default: addressForm.value.is_default ? 1 : 0,
  }
  try {
    if (editingAddressID.value) {
      await put(`/api/v1/addresses/${editingAddressID.value}`, payload)
    } else {
      await post('/api/v1/addresses', payload)
    }
    await refreshAddresses()
    closeAddressEditor()
  } catch (error: any) {
    alert(error?.message || t('userCenter.saveFailed'))
  } finally {
    savingAddress.value = false
  }
}

async function removeAddress(addr: any) {
  const id = Number(addr?.id || 0)
  if (!id) return
  if (!window.confirm(t('userCenter.deleteConfirm'))) return
  try {
    await del(`/api/v1/addresses/${id}`)
    await refreshAddresses()
  } catch (error: any) {
    alert(error?.message || t('userCenter.deleteFailed'))
  }
}

async function unfavorite(productID: number) {
  const id = Number(productID || 0)
  if (!id) return
  try {
    await del(`/api/v1/products/${id}/favorite`)
    favorites.value = favorites.value.filter((row: any) => Number(row.id) !== id)
    favoriteTotal.value = Math.max(0, favoriteTotal.value - 1)
  } catch (error: any) {
    alert(error?.message || t('userCenter.unfavoriteFailed'))
  }
}

function toProductDetail(id: number) {
  router.push(`/product/${id}`)
}
</script>
