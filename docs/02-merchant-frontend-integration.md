# Merchant Frontend → Travel API

## Runtime boundary

```text
merchant-frontend
      │ REST / HTTP
      ▼
  travel-api :9200
      │ gRPC
      ▼
  travel-rpc :9201
      │
      ▼
    Ent / MySQL
```

`merchant-frontend` does not call `travel-rpc`, `merchant-api2`, or `merchant-rpc` for Travel business operations.

## Current Travel endpoints

- `GET /health`
- `GET /api/travel/products`
- `GET /api/travel/products/:id`
- `POST /api/travel/inventory/check`
- `POST /api/travel/orders`
- `GET /api/travel/orders/:orderNo`

## Merchant context

The Travel API currently receives tenant/merchant context from `X-Tenant-ID` and `X-Merchant-ID`, then forwards those values as gRPC metadata to travel-rpc. Customer context is optional.

The Bearer token is retained by the frontend for the merchant session. Token validation must be implemented at the Travel API boundary before production exposure; the current identity headers are a development integration boundary and must not be treated as a public trust mechanism.

## Frontend integration

The Vite development proxy sends `/api/travel/*` to `travel-api`. The frontend Travel client attaches the current merchant session token and tenant/merchant context and never talks directly to gRPC.
