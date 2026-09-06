package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route post /api/travel/merchant/orders/{orderNo}/refund/handle  HandleRefund
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: HandleRefundReq
//
// Responses:
//  200: HandleRefundResp

func HandleRefundHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logic.SetPathParams(r.Context(), r)
		var req types.HandleRefundReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(ctx, w, err)
			return
		}

		l := logic.NewHandleRefundLogic(ctx, svcCtx)
		resp, err := l.HandleRefund(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
		} else {
			httpx.OkJsonCtx(ctx, w, resp)
		}
	}
}
