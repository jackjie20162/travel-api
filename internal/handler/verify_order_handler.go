package handler

import (
	"net/http"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func VerifyOrderHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logic.SetPathParams(r.Context(), r)
		var req types.VerifyOrderReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(ctx, w, err)
			return
		}

		l := logic.NewVerifyOrderLogic(ctx, serverCtx)
		resp, err := l.VerifyOrder(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
			return
		}
		httpx.OkJsonCtx(ctx, w, resp)
	}
}
