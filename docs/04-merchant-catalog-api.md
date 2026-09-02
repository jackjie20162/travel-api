# Merchant Catalog REST API

所有商户商品、套餐、库存管理请求均进入 `travel-api`，再由 `TravelManagementService` 调用 `travel-rpc`。

## 商品

- `POST /api/travel/merchant/products`
- `PUT /api/travel/merchant/products/:id`
- `POST /api/travel/merchant/products/:id/publish`

## 套餐

- `POST /api/travel/merchant/products/:id/packages`
- `GET /api/travel/merchant/products/:id/packages`

## 库存

- `POST /api/travel/merchant/packages/:id/inventory`
- `GET /api/travel/merchant/packages/:id/inventory`

## 作用域

租户和商户 ID 来自 API 网关收到的认证上下文，并通过 gRPC metadata 传给 travel-rpc；业务请求体不允许自行指定 tenant/merchant。

当前本地联调仍使用 `X-Tenant-ID`、`X-Merchant-ID` 作为上下文入口。生产环境需要在 API 边界替换为真实 Bearer/JWT 校验，不能信任客户端任意填写的租户/商户头。

## 闭环验收

创建商品 → 创建套餐 → 写入日期/时段/价格/库存 → 发布 → 消费端读取 → 预占库存 → 创建订单。
