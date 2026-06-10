# PC 装修多版本设计

## 1. 背景与目标

现有装修模块已经在后端与移动端首页编辑器中支持多副本管理：同一 `page_key` 下通过 `variant_key` 区分不同装修版本，并采用单活发布规则。PC 首页装修当前仍按单版本方式调用 `/decor/pc`、`/decor/pc/publish`，导致无法创建、切换、重命名、删除和发布不同 PC 装修副本。

本次目标是在不新增接口路径、不新增数据库字段的前提下，让 PC 装修使用与移动端首页一致的多副本能力：

- 后台 PC 装修支持版本下拉、复制版本、重命名版本、删除未发布版本。
- 保存草稿与发布上线按当前 `variant_key` 生效。
- 发布任意 PC 副本时，该 `page_key=pc` 下其他副本自动下线，保持单活发布。
- Mock 环境与真实接口行为一致，便于演示和前端联调。
- 文档直接描述当前最新装修接口与 PC 多版本能力。

## 2. 设计原则

1. **复用现有接口**：继续使用 `GET/PUT/POST /admin/api/decor/:page_key...` 通用接口，PC 仅传入 `page_key=pc`。
2. **不拆新模型**：后端 `DecorPage`、`variant_key`、`variant_name`、`is_published` 结构已经满足 PC 多版本，不新增 PC 专用表或接口。
3. **编辑器行为一致**：PC 装修顶部版本操作与移动端装修保持同一交互语义，减少后台学习成本。
4. **载荷结构不回退**：PC 版本内容仍使用当前 `PcDecorPagePayload = { pageStyle, components }`，每个副本保存完整页面配置。
5. **Mock 同步**：Admin Mock 的 PC 装修也维护数组化副本源，覆盖 variants、读取、保存、发布、复制、重命名和删除。

## 3. 方案对比与决策

### 方案 A：复用通用多副本接口（采用）

- PC 编辑器调用 `/decor/pc/variants`、`/decor/pc?variant=...`、`/decor/pc/publish?variant=...` 等现有接口。
- 后端不变，仅补测试确认 `page_key=pc` 的副本能力。
- 成本低，符合现有架构与接口约束。

### 方案 B：为 PC 新增专用多版本接口

- 新增 `/decor/pc/versions` 等专用接口。
- 表达直接，但会重复已有 `variant` 能力，增加维护面。
- 不采用。

### 方案 C：前端本地保存多个 PC 版本

- 不改后端或 Mock，只在 PC payload 内嵌 `versions`。
- 会破坏单活发布、后台列表与前台读取的通用模型。
- 不采用。

决策：采用方案 A。

## 4. 架构设计

### 4.1 后端

后端现有通用接口已经支持任意 `page_key`：

- `GET /admin/api/decor/:page_key/variants`
- `GET /admin/api/decor/:page_key?variant=<variant_key>`
- `PUT /admin/api/decor/:page_key?variant=<variant_key>`
- `POST /admin/api/decor/:page_key/publish?variant=<variant_key>`
- `POST /admin/api/decor/:page_key/copies`
- `PUT /admin/api/decor/:page_key/variants/:variant_key`
- `DELETE /admin/api/decor/:page_key/variants/:variant_key`

PC 只需要传入 `page_key=pc`。服务层已有 `defaultComponentsForPage("pc")` 返回 PC 默认页面载荷，因此新副本或空副本能拿到 `{ pageStyle, components }` 结构。

后端补充服务测试，覆盖：

- `ListVariants(ctx, merchantID, "pc")` 空数据时返回 default 副本，且 components 是 PC 页面对象。
- `CreateVariantCopy(ctx, merchantID, "pc", ...)` 复制 PC payload，不退化为数组。
- `PublishPage(ctx, merchantID, "pc", variant)` 只影响 `page_key=pc` 下的发布态，不影响 `index`。

### 4.2 Admin PC 编辑器

`PcDecorEditor.vue` 复用移动端 `DecorEditor.vue` 的版本管理状态：

- `variants`
- `currentVariantKey`
- `loadVariants`
- `loadCurrentVariant`
- `copyVariant`
- `renameVariant`
- `deleteVariant`
- `publish`

PC 保存与发布改为：

- 保存：`PUT /decor/pc?variant=${currentVariantKey}`，body 为 `{ components: pagePayload.value }`
- 发布：先保存，再 `POST /decor/pc/publish?variant=${currentVariantKey}`
- 切换版本：`GET /decor/pc?variant=${currentVariantKey}` 后调用现有 `normalizePayload`

PC 顶部工具栏增加版本下拉与版本操作按钮，文案沿用现有 i18n key：`decor.published`、`decor.copyVariant`、`decor.rename`、`decor.deleteVariant`、`decor.saveDraft`、`decor.publish`。

### 4.3 Admin Mock

当前 `pcDecorSource` 是单对象。改为 `pcDecorVariantsSource` 数组，字段与 `decorVariantsSource` 对齐：

- `id`
- `page_key: "pc"`
- `variant_key`
- `variant_name`
- `components`
- `is_published`
- `published_at`

Mock 补齐以下 PC 分支：

- `GET /admin/api/decor/pc/variants`
- `GET /admin/api/decor/pc?variant=...`
- `PUT /admin/api/decor/pc?variant=...`
- `POST /admin/api/decor/pc/publish?variant=...`
- `POST /admin/api/decor/pc/copies`
- `PUT /admin/api/decor/pc/variants/:variant_key`
- `DELETE /admin/api/decor/pc/variants/:variant_key`

发布逻辑仅对 `pcDecorVariantsSource` 单活，不影响移动端 `decorVariantsSource`。

### 4.4 前台读取

Web 前台继续调用 `GET /api/v1/decor/pc` 获取当前已发布 PC 副本。后端 `GetPublishedPage` 已按 `is_published` 返回同一 `page_key` 最新发布副本，无需前台增加版本参数。

预览或指定版本读取可以使用通用前台接口 `GET /api/v1/decor/pc?variant=<variant_key>`，但本次后台 PC 编辑器不依赖该能力。

## 5. 数据流

1. 进入 PC 装修页。
2. Admin 请求 `/decor/pc/variants`。
3. 若当前版本不存在，选择已发布版本；否则选择第一条或 default。
4. Admin 请求 `/decor/pc?variant=<current>` 并归一化为 `PcDecorPagePayload`。
5. 用户编辑页面样式和组件。
6. 保存草稿时写入当前副本。
7. 发布时先保存当前副本，再将当前副本置为已发布，并下线其他 PC 副本。
8. Web 首页读取 `/api/v1/decor/pc`，展示当前已发布 PC 副本。

## 6. 错误处理与边界

- `variant_key=default` 不能删除。
- 已发布副本不能删除；需要先发布其他副本。
- 复制版本时 `new_variant_key` 为空或与来源相同时，由后端返回错误；Mock 维持同等约束。
- 版本载荷解析失败时，PC 编辑器回退 `createDefaultPcDecorPayload()`，避免白屏。
- `currentVariantKey` 不存在时，优先切到已发布副本，其次第一条副本，最后 default。

## 7. 测试与验收

### 7.1 自动化检查

- `server`: 运行 `go test ./plugins/decor/service -v`，覆盖 PC 默认副本、复制、发布单活。
- `admin`: 运行 `npm run build`，覆盖 PC 编辑器与 Mock 类型/语法。
- 如触及文档站构建配置，运行 `npm run build`（工作目录 `docs-site`）。

### 7.2 手工验收

- 打开后台 PC 首页装修，默认副本可正常加载。
- 复制一个 PC 副本，切换后编辑页面背景或组件内容，保存后刷新仍能回显。
- 发布新 PC 副本后，版本下拉只有该副本显示已发布。
- 删除未发布 PC 副本成功；删除 default 或已发布副本被阻止。
- Web 首页展示当前已发布 PC 副本。
- 移动端首页装修版本管理不受 PC 版本操作影响。

## 8. docs-site 影响

本次属于系统功能变更，需要同步更新 `docs-site/docs/api/decor.md`：

- 功能说明改为装修模块支持任意页面多副本，PC 首页装修同样支持。
- 接口说明中明确 `page_key` 可取 `index`、`pc` 等页面标识。
- PC 装修样式系统章节补充版本管理语义：每个版本保存完整 `PcDecorPagePayload`，发布按 `page_key=pc` 单活。
- 部署影响说明保持不新增环境变量、不新增数据库字段；依赖既有装修副本字段与插件迁移。
