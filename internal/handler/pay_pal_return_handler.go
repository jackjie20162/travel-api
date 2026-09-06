package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
)

// swagger:route get /api/travel/payments/paypal/return  PayPalReturn
//

//

//
// Responses:
//  200: PaymentResp

func PayPalReturnHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewPayPalReturnLogic(r.Context(), svcCtx)
		resp, err := l.PayPalReturn()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
