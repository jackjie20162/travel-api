package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
)

// swagger:route get /api/travel/merchant/products/{id}/packages  ListMerchantPackages
//

//

//
// Responses:
//  200: PackageListResp

func ListMerchantPackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logic.SetPathParams(r.Context(), r)
			l := logic.NewListMerchantPackagesLogic(ctx, svcCtx)
		resp, err := l.ListMerchantPackages()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
