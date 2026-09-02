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
- Travel API wiring is present, but API build/test verification remains a separate milestone.

### Next

1. Verify the latest `travel-rpc` build/test workflow.
2. Verify `travel-api` generation/build/tests independently.
3. Harden reservation/order consistency inside `travel-rpc` without introducing cross-service dependencies.
4. Add payment as an internal Travel payment abstraction after the core service is stable.
5. Publish the REST contract as the integration boundary for Web/App/other adapters.
