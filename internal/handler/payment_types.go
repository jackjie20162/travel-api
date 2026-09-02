package handler

type CreatePaymentReq struct { OrderNo string `json:"orderNo"`; Provider string `json:"provider"`; IdempotencyKey string `json:"idempotencyKey"` }
type PaymentResp struct { Id int64 `json:"id"`; PaymentNo string `json:"paymentNo"`; OrderNo string `json:"orderNo"`; Provider string `json:"provider"`; ProviderPaymentId string `json:"providerPaymentId"`; Amount int64 `json:"amount"`; Currency string `json:"currency"`; Status string `json:"status"`; CheckoutURL string `json:"checkoutUrl,omitempty"` }
