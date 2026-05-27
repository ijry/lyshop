# PC 装修高级样式系统设计（无兼容包袱版）

## 1. 背景与目标

当前 PC 装修完成后风格同质化明显，差异主要来自组件内容而非视觉系统。  
本次目标是在 **不新增接口路径** 的前提下，升级 PC 装修数据模型与渲染能力，支持：

- 页面级背景：纯色 / 渐变 / 背景图（含定位、平铺、尺寸、滚动策略）
- 页面级视觉基线：内容宽度、横向留白、模块间距、圆角、阴影
- 组件级样式覆写：边距、内边距、背景、边框、圆角、阴影
- Admin 编辑器即时预览 + Web 端 Mock 同步展示

范围明确为：**首期仅 PC 装修**，但 Mock 结构在 admin/web/app 三端保持一致。

---

## 2. 设计原则

1. **结构清晰优先**：采用 pageStyle + components 分层，避免样式散落在各组件 props。
2. **接口稳定优先**：不新增 `/decor/pc` 相关接口，仅升级载荷结构。
3. **前后端一致优先**：Admin 预览与 Web 前台共用同一语义字段，减少“编辑可见、前台失真”。
4. **Mock 一致优先**：admin/web/app 的 `pcDecor` 使用同一 Schema，避免后续分叉。

---

## 3. 方案对比与决策

### 方案 A（旧兼容升级）
- 在旧 `components` 外补充字段，兼容历史结构。
- 优点：改造最小；缺点：结构长期会继续膨胀。

### 方案 B（本次采用：无兼容包袱的干净 Schema）
- 统一升级为 `PcDecorPage = { pageStyle, components }`。
- 优点：语义清晰、后续扩展稳定、便于主题模板化。
- 成本：首期改动比兼容方案略大，但总成本更低。

### 方案 C（完全自由样式 JSON）
- 灵活但不可控，易出现冲突与维护问题。

**决策：采用方案 B。**

---

## 4. 数据模型设计

```ts
interface PcDecorPage {
  pageStyle: {
    background: {
      mode: 'solid' | 'gradient' | 'image'
      solidColor?: string
      gradient?: {
        angle: number
        stops: Array<{ color: string; position: number }>
      }
      image?: {
        url: string
        size: 'cover' | 'contain' | 'auto' | 'custom'
        customSize?: string
        position: string
        repeat: 'no-repeat' | 'repeat' | 'repeat-x' | 'repeat-y'
        attachment: 'scroll' | 'fixed'
      }
      overlay?: {
        enabled: boolean
        color: string
        opacity: number
      }
    }
    content: {
      maxWidth: number
      gutterX: number
      sectionGap: number
    }
    surface: {
      radius: number
      shadow: 'none' | 'sm' | 'md' | 'lg'
    }
  }
  components: Array<{
    id: string
    type: string
    props: Record<string, any>
    style?: {
      marginTop?: number
      marginBottom?: number
      paddingX?: number
      paddingY?: number
      backgroundColor?: string
      borderRadius?: number
      borderColor?: string
      borderWidth?: number
      shadow?: 'none' | 'sm' | 'md' | 'lg'
    }
  }>
}
```

说明：
- 不再做旧结构兜底（项目未上生产、无历史包袱）。
- 页面级控制“底层风格”，组件级做局部差异化覆写。

---

## 5. 渲染规则

### 5.1 层级顺序
1. 页面背景（`pageStyle.background`）
2. 可选遮罩（`overlay.enabled=true`）
3. 组件内容层（`components`）

### 5.2 继承与覆写
- 全局基线：`content.maxWidth/gutterX/sectionGap` + `surface.radius/shadow`
- 组件若存在 `style` 字段，则按“组件优先”覆写全局
- 组件间距优先用 `marginTop/marginBottom`，未设则回退 `sectionGap`

### 5.3 背景策略
- `solid`：直接使用 `solidColor`
- `gradient`：`linear-gradient(angle, stops...)`
- `image`：由 `url + size + customSize + position + repeat + attachment` 组合

---

## 6. Admin 编辑器设计（PC 装修）

### 6.1 页面样式面板（新增）
- 背景模式切换：纯色 / 渐变 / 图片
- 图片来源：**URL 输入 + 现有上传组件**
- 背景参数：定位、平铺、尺寸、滚动策略
- 遮罩设置：启用、颜色、透明度
- 内容与基线：maxWidth、gutterX、sectionGap、默认圆角、默认阴影

### 6.2 组件样式覆盖（新增到所有 PC 组件编辑器）
- marginTop / marginBottom
- paddingX / paddingY
- backgroundColor
- borderWidth / borderColor
- borderRadius
- shadow

### 6.3 即时预览
- 所有改动实时同步到中间 `PcDecorPreview`
- 不依赖保存接口，保证编辑体验流畅

---

## 7. 接口与存储策略

不新增接口，仅升级数据体语义：

- `GET /admin/api/decor/pc`
- `PUT /admin/api/decor/pc`
- `POST /admin/api/decor/pc/publish`
- `GET /api/v1/pc/decor`

后端存储保持现状（JSON 字段存页面配置对象），避免迁移数据库表结构。

---

## 8. Mock 同步设计（admin/web/app）

### 8.1 admin mock
- `pcDecorSource` 改为 `PcDecorPage` 结构
- `PUT /admin/api/decor/pc` 按完整对象覆盖写入

### 8.2 web mock
- `preset.pcDecor` 改为 `PcDecorPage`
- `GET /api/v1/pc/decor` 返回完整结构

### 8.3 app mock（类型同步）
- `app/mock/presets/types.ts` 中 `pcDecor` 类型升级为 `PcDecorPage`
- 即使首期不消费 PC 渲染，也保证三端数据结构一致

---

## 9. 校验与容错

### 9.1 保存前校验
- `gradient.stops` 至少 2 个
- `position` 范围 `0~100`
- `mode=image` 时 `image.url` 必填
- 数值类字段最小值校验，禁止负值

### 9.2 运行时容错
- 字段缺失回退默认值，禁止白屏
- 无效颜色/尺寸回退默认并在开发环境告警
- 单组件样式异常不影响其他组件

---

## 10. 测试与验收

### 10.1 验收标准
- Admin 中页面样式与组件样式改动可即时预览
- 保存后回显一致
- Web 在 `VITE_MOCK=true` 下按新结构正确展示
- admin/web/app Mock 结构一致且类型检查通过

### 10.2 建议测试
- `PcDecorPreview`：3 种页面背景模式渲染测试
- `DecorRenderer`：全局样式继承 + 组件覆写测试
- Mock 路由断言：`/admin/api/decor/pc` 与 `/api/v1/pc/decor` 返回结构一致

---

## 11. 对 docs-site 的影响

本次属于系统功能变更，实施时需同步更新 `docs-site`，至少覆盖：
- 功能说明（页面级样式 + 组件级覆写）
- 接口/载荷变化（`/decor/pc` 数据结构升级）
- Mock 使用方式与字段说明

