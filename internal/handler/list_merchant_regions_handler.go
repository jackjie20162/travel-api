package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route get /api/travel/merchant/regions  ListMerchantRegions
//
// ── Regions (merchant management) ──
//
// ── Regions (merchant management) ──
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: RegionListReq
//
// Responses:
//  200: RegionListResp

func ListMerchantRegionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegionListReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewListMerchantRegionsLogic(r.Context(), svcCtx)
		resp, err := l.ListMerchantRegions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
