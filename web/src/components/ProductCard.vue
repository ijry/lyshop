<template>
  <div @click="$router.push(`/product/${product.id}`)"
    class="product-card group bg-white rounded-xl border border-gray-100 overflow-hidden cursor-pointer hover:shadow-lg hover:border-gray-200 hover:-translate-y-1 transition-all duration-300">
    <div class="aspect-square overflow-hidden bg-gray-50">
      <img :src="product.cover" :alt="product.title"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
    </div>
    <div class="p-4">
      <h3 class="product-title text-sm font-medium text-gray-800 line-clamp-2 leading-5 mb-2 transition-colors">
        {{ product.title }}
      </h3>
      <p v-if="product.subtitle" class="text-xs text-gray-400 mb-3 line-clamp-1">{{ product.subtitle }}</p>
      <div class="flex-between">
        <div class="flex items-baseline gap-1.5">
          <span class="text-lg font-bold" :style="{ color: 'var(--color-price, #ef4444)' }">¥{{ product.price }}</span>
          <span v-if="product.origin_price" class="text-xs text-gray-300 line-through">¥{{ product.origin_price }}</span>
        </div>
        <span class="text-xs text-gray-400">{{ product.sales }}{{ $t('productCard.peoplePaid') }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ product: any }>()
</script>

<style scoped>
.product-card:hover .product-title {
  color: var(--color-primary, #dc2626);
}
</style>
