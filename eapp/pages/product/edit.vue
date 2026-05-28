<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { createProduct, getProductDetail, updateProduct } from '@/api/product'

const id = ref(0)
const saving = ref(false)
const form = reactive<any>({
  title: '',
  subtitle: '',
  price: 0,
  stock: 0,
  status: 1,
})

async function loadData() {
  if (!id.value) return
  const data: any = await getProductDetail(id.value)
  Object.assign(form, {
    title: data?.title || '',
    subtitle: data?.subtitle || '',
    price: Number(data?.price || 0),
    stock: Number(data?.stock || 0),
    status: Number(data?.status || 0) === 1 ? 1 : 0,
  })
}

async function save() {
  if (!form.title.trim()) {
    uni.showToast({ title: '请输入商品标题', icon: 'none' })
    return
  }
  saving.value = true
  try {
    const payload = {
      product: {
        title: form.title.trim(),
        subtitle: String(form.subtitle || '').trim(),
        price: Number(form.price || 0),
        stock: Number(form.stock || 0),
        status: Number(form.status || 0) === 1 ? 1 : 0,
      },
    }
    if (id.value) {
      await updateProduct(id.value, payload)
    } else {
      await createProduct(payload)
    }
    uni.showToast({ title: '保存成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 350)
  } finally {
    saving.value = false
  }
}

onLoad((opts) => {
  id.value = Number(opts?.id || 0)
  loadData()
})
</script>

<template>
  <view class="page">
    <view class="card">
      <up-form>
        <up-form-item label="标题"><up-input v-model="form.title" placeholder="请输入商品标题" /></up-form-item>
        <up-form-item label="副标题"><up-input v-model="form.subtitle" placeholder="请输入副标题" /></up-form-item>
        <up-form-item label="价格"><up-input v-model="form.price" type="digit" inputmode="decimal" /></up-form-item>
        <up-form-item label="库存"><up-input v-model="form.stock" type="number" inputmode="numeric" /></up-form-item>
      </up-form>
      <view class="mt-24rpx" />
      <up-button type="primary" :loading="saving" @click="save">保存</up-button>
    </view>
  </view>
</template>

<style scoped>
.page { min-height: 100vh; background: var(--eapp-bg); padding: 20rpx; box-sizing: border-box; }
.card { background: #fff; border: 1px solid var(--eapp-border); border-radius: 22rpx; padding: 24rpx; }
</style>
