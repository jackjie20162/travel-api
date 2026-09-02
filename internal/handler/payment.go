package handler

import (
    "errors"
    "net/http"

    "gitee.com/meinongyihe/travel-api/internal/svc"
    "gitee.com/meinongyihe/travel-rpc/travel"
    "github.com/zeromicro/go-zero/rest/httpx"
)

type PaymentHandler struct { svcCtx *svc.ServiceContext }

func NewPaymentHandler(svcCtx *svc.ServiceContext) *PaymentHandler { return &PaymentHandler{svcCtx: svcCtx} }

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreatePaymentReq
    if err := httpx.Parse(r, &req); err != nil { httpx.Error(w, err); return }
    if req.OrderNo == "" || req.Provider == "" || req.IdempotencyKey == "" {
        httpx.Error(w, errors.New("orderNo, provider and idempotencyKey are required")); return
    }
    ctx, err := rpcContext(r)
    if err != nil { httpx.Error(w, err); return }
    p, err := h.svcCtx.PaymentClient.Create(ctx, &travel.CreatePaymentRequest{
        OrderNo: req.OrderNo, Provider: req.Provider, IdempotencyKey: req.IdempotencyKey,
    })
    if err != nil { httpx.Error(w, err); return }
    httpx.OkJson(w, toPaymentResp(p))
}

func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
    paymentNo := r.PathValue("paymentNo")
    if paymentNo == "" { httpx.Error(w, errors.New("invalid payment number")); return }
    ctx, err := rpcContext(r)
    if err != nil { httpx.Error(w, err); return }
    p, err := h.svcCtx.PaymentClient.Get(ctx, &travel.PaymentNoRequest{PaymentNo: paymentNo})
    if err != nil { httpx.Error(w, err); return }
    httpx.OkJson(w, toPaymentResp(p))
}

func toPaymentResp(p *travel.Payment) PaymentResp {
    return PaymentResp{Id:p.Id, PaymentNo:p.PaymentNo, OrderNo:p.OrderNo, Provider:p.Provider, ProviderPaymentId:p.ProviderPaymentId, Amount:p.Amount, Currency:p.Currency, Status:p.Status}
}
