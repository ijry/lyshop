<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import AppIcon from '@/components/common/AppIcon.vue'
import { activityMeta, formatMoney, parseSkuLabel, type MarketingKind } from '@/utils/marketing'

const props = defineProps<{ row: any; kind: MarketingKind }>()
const emit = defineEmits<{ edit: []; remove: [] }>()

// 外链封面在离线 / 小程序白名单未配置时会加载失败，降级为线性占位图
const coverFailed = ref(false)
watch(
  () => props.row?.product_cover,
  () => {
    coverFailed.value = false
  },
)

const meta = computed(() => activityMeta(props.kind))
const mainPrice = computed(() => Number(props.row?.[meta.value.priceField] || 0))
const secondaryPrice = computed(() => {
  if (meta.value.secondaryField) return Number(props.row?.[meta.value.secondaryField] || 0)
  return Number(props.row?.sku_price || 0)
})
const skuLabel = computed(() => parseSkuLabel(props.row?.sku_attrs, Number(props.row?.sku_id || 0)))
const discount = computed(() => {
  const original = Number(props.row?.sku_price || 0)
  if (props.kind === 'bargain' || !original || !mainPrice.value) return ''
  const value = Math.round((mainPrice.value / original) * 100) / 10
  return value > 0 && value < 10 ? `${value} 折` : ''
})
</script>

<template>
  <view class="product-card">
    <view class="product-card__main" hover-class="eapp-hover" @click="emit('edit')">
      <view class="product-card__cover">
        <AppIcon name="product" :size="44" color="faint" />
        <image
          v-if="row.product_cover && !coverFailed"
          class="product-card__cover-img"
          :src="row.product_cover"
          mode="aspectFill"
          @error="coverFailed = true"
        />
      </view>

      <view class="product-card__info">
        <text class="product-card__title">{{ row.product_title || `商品 #${row.product_id}` }}</text>
        <text class="product-card__sku">{{ skuLabel }}</text>

        <view class="product-card__prices">
          <view class="product-card__price-main">
            <text class="product-card__price-label">{{ meta.priceLabel }}</text>
            <text class="product-card__price-value eapp-num">{{ formatMoney(mainPrice) }}</text>
          </view>
          <view v-if="secondaryPrice" class="product-card__price-sub">
            <text class="product-card__price-label">{{ meta.secondaryLabel || '原价' }}</text>
            <text class="product-card__price-old eapp-num">{{ formatMoney(secondaryPrice) }}</text>
          </view>
          <text v-if="discount" class="product-card__discount eapp-num">{{ discount }}</text>
        </view>
      </view>
    </view>

    <view class="product-card__chips">
      <view class="chip">
        <text class="chip__label">限购</text>
        <text class="chip__value eapp-num">{{ row.limit_per_order || '不限' }}</text>
      </view>
      <view class="chip">
        <text class="chip__label">库存</text>
        <text class="chip__value eapp-num">{{ row.total_stock_limit || '不限' }}</text>
      </view>
      <view class="chip chip--sold">
        <text class="chip__label">已售</text>
        <text class="chip__value eapp-num">{{ row.sold_qty || 0 }}</text>
      </view>
    </view>

    <view class="product-card__foot">
      <text class="product-card__sku-no eapp-num">SKU #{{ row.sku_id }}</text>
      <view class="product-card__actions">
        <view class="mini-action" hover-class="eapp-hover" @click="emit('edit')">
          <AppIcon name="edit" :size="24" color="primary" />
          <text>编辑</text>
        </view>
        <view class="mini-action mini-action--danger" hover-class="eapp-hover" @click="emit('remove')">
          <AppIcon name="close" :size="24" color="danger" />
          <text>移除</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.product-card {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  padding: var(--eapp-space-4);
  box-sizing: border-box;
}

.product-card__main {
  display: flex;
  gap: var(--eapp-space-3);
}

.product-card__cover {
  position: relative;
  width: 136rpx;
  height: 136rpx;
  border-radius: var(--eapp-radius);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--eapp-surface-subtle);
  border: 1rpx solid var(--eapp-border);
  overflow: hidden;
  box-sizing: border-box;
}

.product-card__cover-img {
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
}

.product-card__info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.product-card__title {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-card__sku {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-card__prices {
  display: flex;
  align-items: baseline;
  gap: var(--eapp-space-3);
  flex-wrap: wrap;
  margin-top: var(--eapp-space-2);
}

.product-card__price-main,
.product-card__price-sub {
  display: flex;
  align-items: baseline;
  gap: 6rpx;
}

.product-card__price-label {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.product-card__price-value {
  font-size: var(--eapp-fs-md);
  font-weight: 700;
  color: var(--eapp-primary);
}

.product-card__price-old {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-faint);
  text-decoration: line-through;
}

.product-card__discount {
  font-size: var(--eapp-fs-xs);
  font-weight: 600;
  color: var(--eapp-danger);
  background: var(--eapp-danger-soft);
  border-radius: var(--eapp-radius-sm);
  padding: 2rpx 10rpx;
}

.product-card__chips {
  display: flex;
  gap: var(--eapp-space-2);
  margin-top: var(--eapp-space-4);
}

.chip {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6rpx;
  padding: var(--eapp-space-2) 4rpx;
  border-radius: var(--eapp-radius-sm);
  background: var(--eapp-surface-subtle);
  border: 1rpx solid var(--eapp-divider);
  box-sizing: border-box;
}

.chip--sold {
  background: var(--eapp-warning-soft);
  border-color: var(--eapp-warning-soft);
}

.chip__label {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.chip__value {
  font-size: var(--eapp-fs-sm);
  font-weight: 600;
  color: var(--eapp-text-secondary);
}

.chip--sold .chip__value {
  color: var(--eapp-warning);
}

.product-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--eapp-space-3);
  margin-top: var(--eapp-space-3);
  padding-top: var(--eapp-space-3);
  border-top: 1rpx solid var(--eapp-divider);
}

.product-card__sku-no {
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-faint);
}

.product-card__actions {
  display: flex;
  gap: var(--eapp-space-2);
}

.mini-action {
  display: flex;
  align-items: center;
  gap: 4rpx;
  height: 54rpx;
  padding: 0 var(--eapp-space-3);
  border-radius: var(--eapp-radius-sm);
  border: 1rpx solid var(--eapp-border);
  background: var(--eapp-surface-subtle);
  font-size: var(--eapp-fs-xs);
  color: var(--eapp-text-secondary);
  box-sizing: border-box;
}

.mini-action--danger {
  color: var(--eapp-danger);
  border-color: var(--eapp-danger-soft);
  background: var(--eapp-danger-soft);
}
</style>
