package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
)

// swagger:route post /api/travel/my/orders/{orderNo}/cancel  CancelOrder
//

//

//
// Responses:
//  200: CancelOrderResp

func CancelOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logic.SetPathParams(r.Context(), r)
		l := logic.NewCancelOrderLogic(ctx, svcCtx)
		resp, err := l.CancelOrder()
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
		} else {
			httpx.OkJsonCtx(ctx, w, resp)
		}
	}
}
