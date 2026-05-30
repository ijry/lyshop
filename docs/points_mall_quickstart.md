# 积分商城插件快速启动指南

## 1. 启用插件

### 创建配置文件
```bash
cd server
cp config.example.yaml config.yaml
```

### 编辑配置文件
修改 `config.yaml` 中的数据库连接信息：
```yaml
database:
  dsn: "root:your_password@tcp(127.0.0.1:3306)/lyshop?charset=utf8mb4&parseTime=True&loc=Local"
```

确认 `points_mall` 插件已在启用列表中：
```yaml
plugins:
  enabled:
    - product
    - order
    - marketing
    - points_mall  # ✅ 积分商城插件
    - vip
    - checkin
```

## 2. 启动后端服务

```bash
cd server
go mod tidy
go run main.go -config config.yaml
```

服务启动后会自动：
- 注册 points_mall 插件
- 创建数据库表（points_logs, points_products, points_exchanges, points_configs）
- 注册 API 路由

## 3. 启动前端应用

### Admin 端（管理后台）
```bash
cd admin
npm install  # 首次运行
npm run dev
```
访问：http://localhost:5173

### Eapp 端（商家端）
```bash
cd eapp
npm install  # 首次运行
npm run dev:h5
```

### App 端（用户端）
```bash
cd app
npm install  # 首次运行
npm run dev:h5
```

### Web 端（用户 Web）
```bash
cd web
npm install  # 首次运行
npm run dev
```

## 4. 功能测试流程

### 4.1 Admin 端测试

1. **登录管理后台**
   - 访问 http://localhost:5173
   - 使用管理员账号登录

2. **创建积分商品**
   - 导航到：积分商城 → 积分商品
   - 点击"新增商品"
   - 填写商品信息：
     - 标题：测试商品
     - 类型：实物/优惠券/虚拟
     - 积分价格：1000
     - 库存：100
     - 封面图片URL
     - 商品描述
   - 保存

3. **查看积分统计**
   - 导航到：积分商城 → 积分统计
   - 查看累计发放、消耗、余额等数据

4. **管理兑换记录**
   - 导航到：积分商城 → 兑换记录
   - 查看用户兑换记录
   - 对实物商品进行发货操作

5. **查看积分日志**
   - 导航到：积分商城 → 积分日志
   - 查看所有积分变动记录
   - 可以手动调整用户积分

### 4.2 App 端测试

1. **浏览积分商城**
   - 访问积分商城首页
   - 查看我的积分余额
   - 浏览商品列表

2. **兑换商品**
   - 点击商品查看详情
   - 点击"立即兑换"
   - 确认兑换

3. **查看兑换记录**
   - 进入"我的兑换记录"
   - 查看兑换状态
   - 对已发货商品确认收货

4. **签到赚积分**
   - 进入签到页面
   - 每日签到获得积分
   - 查看积分日志

### 4.3 集成功能测试

1. **签到赠送积分**
   ```
   用户签到 → 自动调用 points_mall 服务 → 增加积分 → 记录日志
   ```

2. **订单抵扣积分**
   ```
   下单时选择使用积分 → 积分计算器计算抵扣金额 → 扣减积分
   ```

3. **兑换优惠券**
   ```
   兑换优惠券类商品 → 自动发放优惠券到用户账户 → 可在订单中使用
   ```

4. **订单完成赠送积分**（需要在 order 插件中集成）
   ```go
   // 在订单完成时调用
   import pmservice "github.com/ijry/lyshop/server/plugins/points_mall/service"
   
   if newStatus == ordermodel.OrderStatusCompleted {
       pmservice.GrantOrderPoints(ctx, orderID)
   }
   ```

## 5. API 测试

### 使用 curl 测试

```bash
# 获取积分商品列表
curl http://localhost:8080/api/v1/points/products

# 获取用户积分余额（需要 token）
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:8080/api/v1/points/balance

# 兑换商品（需要 token）
curl -X POST \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"qty":1}' \
     http://localhost:8080/api/v1/points/products/1/exchange

# 管理端：获取积分统计（需要 admin token）
curl -H "Authorization: Bearer ADMIN_TOKEN" \
     http://localhost:8080/admin/api/points/stats
```

## 6. 数据库检查

连接数据库查看表结构：

```sql
-- 查看积分日志
SELECT * FROM points_logs ORDER BY id DESC LIMIT 10;

-- 查看积分商品
SELECT * FROM points_products;

-- 查看兑换记录
SELECT * FROM points_exchanges ORDER BY id DESC LIMIT 10;

-- 查看用户积分余额
SELECT id, username, points FROM users WHERE points > 0;
```

## 7. 常见问题

### Q1: 插件启动失败
**检查**：
- 数据库连接是否正确
- `points_mall` 是否在 `config.yaml` 的 `plugins.enabled` 列表中
- 依赖的 `marketing` 插件是否已启用

### Q2: 前端页面空白
**检查**：
- 后端服务是否正常运行
- API 路由是否正确注册
- 浏览器控制台是否有错误

### Q3: 兑换失败
**检查**：
- 用户积分是否足够
- 商品库存是否充足
- 商品是否已上架（status=1）

### Q4: 积分未增加
**检查**：
- 积分日志表是否有记录
- 用户表的 points 字段是否更新
- 是否有事务回滚

## 8. 性能优化建议

1. **数据库索引**
   - points_logs: user_id, type, created_at
   - points_products: status, type, sort
   - points_exchanges: user_id, status, created_at

2. **缓存策略**
   - 积分商品列表可以缓存 5 分钟
   - 用户积分余额可以缓存 1 分钟
   - 积分统计数据可以缓存 10 分钟

3. **并发控制**
   - 兑换商品时使用数据库行锁
   - 积分扣减使用事务保证原子性

## 9. 监控指标

建议监控以下指标：
- 每日积分发放总量
- 每日积分消耗总量
- 商品兑换成功率
- 兑换订单平均处理时间
- 积分余额分布

## 10. 下一步扩展

- [ ] 实现积分过期机制
- [ ] 添加积分排行榜
- [ ] 支持积分转赠
- [ ] 添加积分任务系统
- [ ] 实现积分抽奖功能
- [ ] 支持积分等级制度

---

## 技术支持

如有问题，请查看：
- 完整文档：`docs/points_mall_implementation.md`
- 插件代码：`server/plugins/points_mall/`
- 前端代码：`admin/src/views/points-mall/`
