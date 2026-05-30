# 营销活动商品选择器弹窗设计

**日期：** 2026-05-30  
**范围：** admin 后台管理前端  
**涉及文件：** `ActivityProductManageBase.vue`、新增 `ProductSelectorModal.vue`

---

## 背景与目标

当前营销活动（秒杀、拼团、砍价）在添加商品时，使用 `<select>` 下拉框展示所有商品（一次性加载 200 条），存在以下问题：

- 商品数量多时下拉框难以查找
- 无法展示商品封面图等视觉信息
- 无搜索/分页能力，体验差

**目标：** 将商品选择方式改为弹窗，支持关键字搜索、分类筛选、分页浏览，并展示商品完整信息。

---

## 方案选择

采用**方案 A：创建独立的 `ProductSelectorModal.vue` 组件**。

理由：
- 不影响现有功能，风险最低
- 组件职责清晰，可在所有营销活动类型中复用
- 不需要改造现有 ProductList 页面

---

## 组件设计

### 新增：`admin/src/components/product/ProductSelectorModal.vue`

**Props：**

```typescript
interface Props {
  modelValue: boolean     // 控制显示/隐藏（v-model）
  excludeIds?: number[]   // 已选商品 ID 列表，这些商品在列表中置灰不可选
}
```

**Emits：**

```typescript
interface Emits {
  'update:modelValue': (value: boolean) => void  // 关闭弹窗
  'select': (product: Product) => void           // 用户选中某商品
}
```

**内部 Product 类型：**

```typescript
interface Product {
  id: number
  title: string
  cover: string
  price: number
  stock: number
  favorite_count: number
  status: number  // 1=上架, 0=下架
}
```

**内部状态：**

```typescript
const query = ref({ keyword: '', category_id: '', page: 1, size: 20 })
const products = ref<Product[]>([])
const categories = ref<any[]>([])
const total = ref(0)
```

**数据加载：**
- 使用现有的 `getProducts(query)` 和 `getCategories()` API（来自 `@/api/plugins`）
- 弹窗打开时自动加载第一页数据
- 弹窗关闭时重置搜索条件和分页

---

## UI 布局

```
┌─────────────────────────────────────────────────────────┐
│  选择商品                                          [✕]   │
├─────────────────────────────────────────────────────────┤
│  [关键字搜索...] [分类下拉] [搜索按钮]                    │
├─────────────────────────────────────────────────────────┤
│  ID │ 商品（封面+标题）│ 价格 │ 库存 │ 收藏数 │ 状态     │
│─────┼──────────────────┼──────┼──────┼────────┼────────│
│  1  │ [img] 商品标题   │ ¥99  │ 100  │  20    │ 上架   │
│  2  │ [img] 商品标题   │ ¥199 │  50  │   5    │ 下架   │
│  ...                                                    │
├─────────────────────────────────────────────────────────┤
│  共 N 条          [上一页]  [下一页]                     │
└─────────────────────────────────────────────────────────┘
```

**样式规范：**
- 模态框宽度：`w-[1120px] max-w-[96vw]`，最大高度 `max-h-[92vh]`，内容区可滚动
- 遮罩：`fixed inset-0 bg-black/35 flex items-center justify-center z-50`
- 点击遮罩关闭弹窗
- 已在 `excludeIds` 中的商品行：`opacity-40 cursor-not-allowed`，不可点击
- 可点击商品行：`hover:bg-blue-50 cursor-pointer`，点击整行触发选择
- 选中后自动关闭弹窗

---

## 修改：`ActivityProductManageBase.vue`

**变更点：**

1. 移除 `loadProductOptions()` 函数（不再需要预加载 200 条商品）
2. 移除 `productOptions` ref
3. 将表单中的商品 `<select>` 替换为：
   - 一个显示已选商品的只读文本框（或空状态提示）
   - 一个"选择商品"按钮，点击打开 `ProductSelectorModal`
4. 添加 `showProductSelector` ref 控制弹窗显示
5. 在 `onSelectProduct` 回调中接收选中的商品对象，填充 `selectedProductID` 并触发 SKU 加载

**筛选栏的商品下拉框（filterProductID）：**
- 保留现有的 `<select>` 筛选框，但改为使用弹窗选择器
- 或者简化为：移除商品筛选下拉框，仅保留关键字搜索（因为弹窗已经提供了更好的搜索体验）
- 决策：**移除** `filterProductID` 筛选下拉框，减少冗余，用关键字搜索代替

---

## 交互流程

```
用户点击"新增"按钮
  → 打开活动商品表单（右侧抽屉）
  → 点击"选择商品"按钮
    → 打开 ProductSelectorModal
    → 用户可搜索/翻页浏览商品
    → 点击某商品行
      → 弹窗关闭
      → 表单中显示已选商品名称
      → 自动加载该商品的 SKU 列表
  → 用户选择 SKU
  → 填写价格、库存限制等
  → 保存
```

---

## 不在本次范围内

- 多选模式（当前只需单选）
- 商品分类树形筛选（使用现有的平铺分类下拉即可）
- 优惠券页面的商品选择（CouponList.vue 目前不涉及商品选择，不需要修改）

---

## 文件变更清单

| 操作 | 文件 |
|------|------|
| 新增 | `admin/src/components/product/ProductSelectorModal.vue` |
| 修改 | `admin/src/views/marketing/ActivityProductManageBase.vue` |
