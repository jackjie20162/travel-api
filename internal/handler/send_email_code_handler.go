package handler

import (
	"net/http"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendEmailCodeHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendEmailCodeReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewSendEmailCodeLogic(r.Context(), serverCtx)
		resp, err := l.SendEmailCode(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
