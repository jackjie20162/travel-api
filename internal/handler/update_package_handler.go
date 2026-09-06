package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route put /api/travel/merchant/packages/{id}  UpdatePackage
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: PackageReq
//
// Responses:
//  200: Package

func UpdatePackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PackageReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		ctx := logic.SetPathParams(r.Context(), r)
			l := logic.NewUpdatePackageLogic(ctx, svcCtx)
		resp, err := l.UpdatePackage(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
