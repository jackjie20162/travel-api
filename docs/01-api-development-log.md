# Travel API Development Log

## 2026-09-02 — Standalone Travel API → Travel RPC

### Architecture decision

Travel is implemented as a **pluggable standalone microservice**. The service boundary is intentionally kept narrow:

```text
Web / H5 / App / external adapters
              │
              ▼ REST/HTTP
          travel-api
              │ gRPC
              ▼
          travel-rpc
              │
              ▼
           MySQL/Ent
```

`travel-api` and `travel-rpc` are the only components in the Travel service runtime. Other business services must not be embedded into this service or become required runtime dependencies.

### Implemented

- `travel-api` owns REST routes, HTTP validation, authentication/context extraction and RPC orchestration.
- `travel-api` creates one gRPC connection to the configured `TravelRpc.Target`.
- `travel-api` owns RPC client construction through `internal/svc/servicecontext.go`.
- `travel-rpc` owns tourism domain rules, inventory, reservations, orders and persistence.
- Tenant/merchant/customer identity is treated as request context; Travel does not depend on another service's internal implementation.
- Product list/detail REST endpoints are wired to `CatalogService`.
- Inventory check REST endpoint is wired to `InventoryService.Check`.
- Order creation REST endpoint is wired to `OrderService.Create`.
- Order query REST endpoint `GET /api/travel/orders/:orderNo` is now registered and wired to `OrderService.Get`.
- Optional `X-Customer-ID` is propagated to `CreateOrderRequest.CustomerId` when present.

### Context boundary

The current API forwards `X-Tenant-ID`, `X-Merchant-ID` and optional `X-Customer-ID` as gRPC metadata. These are transport-level context fields only. Production authentication must establish and verify these values at the Travel API boundary; no downstream service should trust arbitrary public headers without authentication.

### Explicit non-goals

- No direct `travel-api → merchant-api2` dependency.
- No direct `travel-api → merchant-rpc` dependency.
- No dependency on the simple-admin backend at runtime.
- No payment provider SDK inside `travel-api` at this stage.
- No requirement for the total/platform admin service.

Other systems may integrate with Travel later through its stable REST API, but Travel remains deployable and testable by itself.

### Verification status

- RPC protobuf/Ent generation has previously been verified in CI; the latest type-compatibility fixes are awaiting a new workflow result.
- Travel API source wiring has been extended, but a fresh API generation/build/test workflow result is still required.

### Next

1. Verify the latest `travel-rpc` build/test workflow.
2. Verify `travel-api` generation/build/tests independently.
3. Harden reservation/order consistency inside `travel-rpc` without introducing cross-service dependencies.
4. Add payment as an internal Travel payment abstraction after the core service is stable.
5. Publish the REST contract as the integration boundary for Web/App/other adapters.

## 2026-10-17 — Public attraction endpoints & product cover fallback

### Implemented

- `GET /api/travel/attractions/:id` — 公开景点详情（仅 `status=ACTIVE`，按 `X-Merchant-ID` 隔离），直连 ent 查询，不走 RPC。
- `GET /api/travel/attractions/:id/products` — 反向查询：行程节点 `itinerary_stops.type_params.attractionId` 绑定了该景点的 `PUBLISHED` 线路，LIKE 粗筛 + JSON 精校验，按 `min_price` 升序。
- 线路主图兜底 `internal/logic/product_cover_fallback.go`：公开产品列表/详情/景点关联线路三处响应中，`cover_image` 为空的产品自动取行程绑定景点的图（按行程顺序取第一个有图景点，景点优先封面、其次图集第一张）。只补齐响应字段，不回写数据库；商户后台接口不兜底。

### Verification status

- `go build ./...` 通过；两个景点接口在临时实例（9216）上走查通过（列表/详情/关联线路/不存在 id 报错）。
- 主图兜底由 `internal/logic/product_cover_fallback_db_test.go` 验证（只读，需设 `TRAVEL_DB_DSN` 才跑，否则 skip）：4 个 merchant 1 产品清掉主图后，2 个有景点绑定的自动取得景点封面，2 个无绑定的仍为空。
- 景点图集已补齐到每个 6 张（`travel-rpc/docs/sql/dubai_attractions_images_topup.sql`），图片经本地上传服务 9208 自托管，静态回显 GET 200。
