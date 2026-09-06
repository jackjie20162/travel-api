package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/user"
)

// swagger:route post /api/travel/my/orders/{orderNo}/refund  RequestRefund
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: RequestRefundReq
//
// Responses:
//  200: RequestRefundResp

func RequestRefundHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logic.SetPathParams(r.Context(), r)
		// Extract user ID from token
		token := middleware.ExtractToken(r)
		if token != "" {
			u, err := svcCtx.DB.User.Query().Where(user.TokenEQ(token)).Only(r.Context())
			if err == nil {
				ctx = middleware.WithAuthUserID(ctx, int64(u.ID))
			}
		}
		var req types.RequestRefundReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(ctx, w, err)
			return
		}

		l := logic.NewRequestRefundLogic(ctx, svcCtx)
		resp, err := l.RequestRefund(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
		} else {
			httpx.OkJsonCtx(ctx, w, resp)
		}
	}
}
