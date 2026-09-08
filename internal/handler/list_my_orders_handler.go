package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/user"
)

func ListMyOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderListReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// Extract Bearer token → find user → inject userID into context
		token := middleware.ExtractToken(r)
		ctx := r.Context()
		if token != "" {
			u, err := svcCtx.DB.User.Query().Where(user.TokenEQ(token)).Only(r.Context())
			if err != nil {
				logx.Errorf("[ListMyOrders] token lookup failed: %v, token_prefix=%s", err, token[:min(8, len(token))])
			} else {
				ctx = middleware.WithAuthUserID(ctx, int64(u.ID))
			}
		} else {
			logx.Errorf("[ListMyOrders] no token found in request")
		}
		l := logic.NewListMyOrdersLogic(ctx, svcCtx)
		resp, err := l.ListMyOrders(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
