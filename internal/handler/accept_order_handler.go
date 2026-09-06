package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route post /api/travel/merchant/orders/{orderNo}/accept  AcceptOrder
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: AcceptOrderReq
//
// Responses:
//  200: AcceptOrderResp

func AcceptOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logic.SetPathParams(r.Context(), r)
		var req types.AcceptOrderReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(ctx, w, err)
			return
		}

		l := logic.NewAcceptOrderLogic(ctx, svcCtx)
		resp, err := l.AcceptOrder(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
		} else {
			httpx.OkJsonCtx(ctx, w, resp)
		}
	}
}
