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

// swagger:route put /api/travel/user/profile  UpdateProfile
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: UpdateProfileReq
//
// Responses:
//  200: UserProfile

func UpdateProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateProfileReq
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
		l := logic.NewUpdateProfileLogic(ctx, svcCtx)
		resp, err := l.UpdateProfile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
