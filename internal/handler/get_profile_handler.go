package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-rpc/ent/user"
)

// swagger:route get /api/travel/user/profile  GetProfile
//

//

//
// Responses:
//  200: UserProfile

func GetProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := middleware.ExtractToken(r)
		ctx := r.Context()
		if token != "" {
			u, err := svcCtx.DB.User.Query().Where(user.TokenEQ(token)).Only(r.Context())
			if err == nil {
				ctx = middleware.WithAuthUserID(ctx, int64(u.ID))
			}
		}
		l := logic.NewGetProfileLogic(ctx, svcCtx)
		resp, err := l.GetProfile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
