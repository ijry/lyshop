# 接口速查

本页仅整理当前文档与公开示例中可见的 AI 相关接口信息，用于快速检索。

## 文档可见接口

- `GET /admin/api/ai/models`：获取可用模型列表与能力标识。  
- `GET /admin/api/ai/tasks`：查询生成任务列表。  
- `POST /admin/api/ai/generate`：创建生图任务。  
- `GET /admin/api/ai/tasks/:id`：查询单个任务状态与结果。  

## 最小调用流程

1. 先调用 `GET /admin/api/ai/models` 获取可选模型。  
2. 通过 `POST /admin/api/ai/generate` 提交生成请求。  
3. 使用 `GET /admin/api/ai/tasks/:id` 轮询任务结果。  
4. 成功后将返回图片应用到商品素材。  

## 一致性说明

> 接口字段以 docs-site 当前 API 文档与已公开示例为准。  
> 若后端能力先行变更，请先同步 API 文档，再更新本页。
