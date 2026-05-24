# 装修接口

## 说明

装修模块支持页面组件化配置、草稿保存、发布上线，以及首页多副本管理（单活发布）。

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

## 多副本（首页装修）

- 同一 `page_key`（如 `index`）下可存在多个副本（`variant_key`）。
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
- 前台 `GET /api/v1/index/decor` 仅返回当前已发布副本。
- 若无已发布副本，则回退到 `default` 副本。

## 部署与配置影响

- 装修表结构新增副本字段（`variant_key`、`variant_name`、`is_published`）。
- 旧数据会在插件安装阶段迁移为 `default` 副本，并根据历史发布时间或更新时间推导发布态。
- 无新增环境变量；部署后重启服务触发插件迁移即可生效。
