package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
)

// swagger:route post /api/travel/payments/paypal/webhook  PayPalWebhook
//

//

//

func PayPalWebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewPayPalWebhookLogic(r.Context(), svcCtx)
		err := l.PayPalWebhook(r)
		if err != nil {
			logx.Errorf("[PayPal Webhook] error: %v, headers: PAYPAL-AUTH-ALGO=%q, PAYPAL-CERT-URL=%q, PAYPAL-TRANSMISSION-ID=%q, PAYPAL-TRANSMISSION-SIG=%q, PAYPAL-TRANSMISSION-TIME=%q",
				err,
				r.Header.Get("PAYPAL-AUTH-ALGO"),
				r.Header.Get("PAYPAL-CERT-URL"),
				r.Header.Get("PAYPAL-TRANSMISSION-ID"),
				r.Header.Get("PAYPAL-TRANSMISSION-SIG"),
				r.Header.Get("PAYPAL-TRANSMISSION-TIME"),
			)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
