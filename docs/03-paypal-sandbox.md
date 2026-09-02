# Travel API — PayPal Sandbox

## Runtime path

```text
merchant-frontend
      ↓ REST
travel-api :9200
      ↓ gRPC
travel-rpc :9201
      ↓
Ent / MySQL
```

PayPal is an external payment provider used by `travel-api`. `merchant-api2` and `merchant-rpc` are not runtime dependencies of this flow.

## Endpoints

- `POST /api/travel/payments` — creates the Travel payment and PayPal Checkout order, returning `checkoutUrl`.
- `GET /api/travel/payments/:paymentNo` — reads payment state using merchant tenant context.
- `GET /api/travel/payments/paypal/return` — receives PayPal buyer return and captures the PayPal order server-side.
- `GET /api/travel/payments/paypal/cancel` — receives buyer cancellation.
- `POST /api/travel/payments/paypal/webhook` — verifies PayPal webhook signatures and consumes `PAYMENT.CAPTURE.COMPLETED`.

## Sandbox URLs

```text
PayPal API: https://api-m.sandbox.paypal.com
Webhook:    https://travel.code688.com/api/travel/payments/paypal/webhook
Return:     https://travel.code688.com/api/travel/payments/paypal/return
Cancel:     https://travel.code688.com/api/travel/payments/paypal/cancel
```

## Environment

Do not commit credentials.

```text
PAYPAL_ENV=sandbox
PAYPAL_CLIENT_ID=<local secret/config>
PAYPAL_CLIENT_SECRET=<local secret/config>
PAYPAL_WEBHOOK_ID=<registered webhook id>
```

The Client Secret is read from configuration/environment only and is never returned to the frontend.

## Checkout flow

1. Travel API creates the internal payment through `travel-rpc`.
2. Travel API obtains a PayPal OAuth access token server-side.
3. Travel API creates a PayPal Orders v2 order with `intent=CAPTURE`.
4. PayPal approval URL is returned as `checkoutUrl`.
5. The merchant frontend redirects the buyer to PayPal.
6. PayPal redirects the buyer to the return URL with `token` and the signed Travel state.
7. Travel API captures the PayPal order server-side.
8. A completed capture transitions the Travel payment to `PAID` and the Travel order to `CONFIRMED` through `travel-rpc`.
9. PayPal `PAYMENT.CAPTURE.COMPLETED` webhook is independently verified and processed as an idempotent confirmation path.

The browser return is not treated as payment proof by itself; the server-side capture result and verified webhook are authoritative.

## Amounts

Travel amounts are stored as integer minor units. PayPal values are converted to major units for currencies with two decimal places; JPY/KRW are sent as integer values.

## Verification status

Implemented in source: PayPal OAuth, Orders v2 create, approval URL, capture, signed return state, webhook signature verification, webhook completion handling, and frontend checkout button.

Not yet production-verified: live HTTPS deployment at the domain, registered PayPal Webhook ID, real Sandbox end-to-end transaction, CI green status, and persistent webhook-event table for audit/replay tracking.
