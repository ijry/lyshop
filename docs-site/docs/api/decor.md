# 装修接口

## 说明

装修模块支持页面组件化配置、草稿保存、发布上线，以及任意页面多副本管理（单活发布）。移动端首页通常使用 `page_key=index`，PC 首页使用 `page_key=pc`。

## 典型接口

- `GET /api/v1/index/decor`
- `GET /api/v1/decor/:page_key?variant=<variant_key>`
- `GET /admin/api/decor/:page_key/variants`
- `GET /admin/api/decor/:page_key?variant=<variant_key>`
- `PUT /admin/api/decor/:page_key?variant=<variant_key>`
- `POST /admin/api/decor/:page_key/publish?variant=<variant_key>`
- `POST /admin/api/decor/:page_key/copies`
- `PUT /admin/api/decor/:page_key/variants/:variant_key`
- `DELETE /admin/api/decor/:page_key/variants/:variant_key`

## 多副本（任意页面装修）

- 同一 `page_key`（如 `index`、`pc`）下可存在多个副本（`variant_key`）。
- `variant_key=default` 为默认副本，不支持删除。
- 后台保存草稿与发布动作均可按 `variant` 指定目标副本。

### 复制副本

`POST /admin/api/decor/:page_key/copies`

```json
{
  "from_variant_key": "default",
  "new_variant_key": "spring_festival_2027",
  "new_variant_name": "春节版 2027"
}
```

## 发布规则（单活）

- 发布某个副本时，系统会将该 `page_key` 下其他副本置为未发布。
- 前台 `GET /api/v1/index/decor` 返回 `index` 当前已发布副本；`GET /api/v1/decor/pc` 返回 PC 首页当前已发布副本。
- 若无已发布副本，则回退到 `default` 副本。

## 部署与配置影响

- 装修表结构新增副本字段（`variant_key`、`variant_name`、`is_published`）。
- 旧数据会在插件安装阶段迁移为 `default` 副本，并根据历史发布时间或更新时间推导发布态。
- 无新增环境变量；部署后重启服务触发插件迁移即可生效。

## 编辑器交互升级（2026-05-27）

- 功能说明：
  - 后台「移动端装修」左侧组件库支持点击直接添加组件（保留拖拽添加能力）。
  - 后台「移动端装修」中间手机预览支持点击已渲染组件，右侧自动激活对应属性编辑面板。
- 接口变化：
  - 无，本次仅前端编辑器交互增强，仍使用现有装修读写接口。
- 部署或配置影响：
  - 无新增环境变量、无数据库变更、无额外部署步骤。

## PC 装修样式系统（2026-05-28）

- 功能说明：
  - 后台「PC 首页装修」支持页面级样式与组件级样式能力。
  - PC 首页装修支持多副本管理，每个副本保存完整 `PcDecorPage` 载荷。
  - PC 商城首页装修信息统一通过 `GET /api/v1/decor/pc` 读取。
  - 页面级样式支持背景模式（纯色/渐变/背景图）、内容布局（最大宽度/左右留白/模块间距）与默认外观（圆角/阴影）。
  - 组件级样式支持外边距、内边距、背景色、边框、圆角和阴影覆盖，按组件单独生效。
- 接口变化：
  - 使用现有装修读写与发布接口：
  - `GET /admin/api/decor/pc`
  - `PUT /admin/api/decor/pc`
  - `POST /admin/api/decor/pc/publish`
  - `GET /api/v1/decor/pc`
  - `GET /admin/api/decor/pc/variants`
  - `GET /admin/api/decor/pc?variant=<variant_key>`
  - `PUT /admin/api/decor/pc?variant=<variant_key>`
  - `POST /admin/api/decor/pc/publish?variant=<variant_key>`
  - `POST /admin/api/decor/pc/copies`
  - `PUT /admin/api/decor/pc/variants/:variant_key`
  - `DELETE /admin/api/decor/pc/variants/:variant_key`
  - `components` 字段为页面载荷对象：`{ "pageStyle": {...}, "components": [...] }`
  - `pageStyle` 负责页面背景、遮罩、布局与默认表面样式；`components[].style` 负责组件级样式覆盖。
- 部署或配置影响：
  - 无新增环境变量、无新增数据库字段、无迁移脚本变更。
