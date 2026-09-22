package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route post /api/travel/payments/stripe/intent  StripeCreateIntent
//
// ── Stripe Payment ──
//
// ── Stripe Payment ──
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: CreatePaymentReq
//
// Responses:
//  200: StripeIntentResp

func StripeCreateIntentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatePaymentReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewStripeCreateIntentLogic(r.Context(), svcCtx)
		resp, err := l.StripeCreateIntent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
