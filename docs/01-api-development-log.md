# Travel API Development Log

## 2026-09-02 — API contract correction after RPC generation verification

### Implemented
- Corrected `go.mod` so travel-api references the dedicated `travel-rpc` repository rather than the obsolete `travel-app` path.
- Added `TravelRpc.Target` to the runtime configuration and defaulted it to `127.0.0.1:9201`.
- Updated `desc/travel.api` from v0.1.1 to v0.1.2.
- Added request validation for product IDs, package IDs, service dates and quantities.
- Restricted product-list page size to 1–100 when supplied.
- Made `reservationKey` required for order creation to preserve idempotent inventory/order semantics.
- Kept tenant, merchant, final amount and currency out of client-authoritative order input.

### Verification
- The latest travel-rpc GitHub Actions run successfully completed protobuf generation and Ent generation.
- The run currently fails at `go test ./...` because the repository requires `go.mod` updates (`go mod tidy`); this is a dependency hygiene failure, not an Ent-generation failure.
- travel-api runtime RPC handler wiring is not yet marked complete because generated RPC client code is not yet committed/consumed by this repository.

### Next
1. Normalize travel-rpc module/dependency files and obtain a green RPC CI run.
2. Generate the travel-api Go code from `desc/travel.api` with goctl.
3. Wire Catalog/Inventory/Order handlers to travel-rpc.
4. Add tenant/auth context propagation and API-level tests.
