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

// swagger:route put /api/travel/user/password  ChangePassword
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ChangePasswordReq
//
// Responses:
//  200: OkResp

func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		token := middleware.ExtractToken(r)
		ctx := r.Context()
		if token != "" {
			u, err := svcCtx.DB.User.Query().Where(user.TokenEQ(token)).Only(r.Context())
			if err == nil {
				ctx = middleware.WithAuthUserID(ctx, int64(u.ID))
			}
		}
		l := logic.NewChangePasswordLogic(ctx, svcCtx)
		resp, err := l.ChangePassword(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
