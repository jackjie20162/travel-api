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

// swagger:route post /api/travel/orders  CreateOrder
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: CreateOrderReq
//
// Responses:
//  200: CreateOrderResp

func CreateOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOrderReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// Extract Bearer token → find user → inject userID as customerID
		ctx := r.Context()
		token := middleware.ExtractToken(r)
		if token != "" {
			u, err := svcCtx.DB.User.Query().Where(user.TokenEQ(token)).Only(r.Context())
			if err == nil {
				ctx = middleware.WithAuthUserID(ctx, int64(u.ID))
			}
		}
		l := logic.NewCreateOrderLogic(ctx, svcCtx)
		resp, err := l.CreateOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
