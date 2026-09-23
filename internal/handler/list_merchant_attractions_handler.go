package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route get /api/travel/merchant/attractions  ListMerchantAttractions
//
// ── Attractions (merchant management) ──
//
// ── Attractions (merchant management) ──
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: AttractionListReq
//
// Responses:
//  200: AttractionListResp

func ListMerchantAttractionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AttractionListReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewListMerchantAttractionsLogic(r.Context(), svcCtx)
		resp, err := l.ListMerchantAttractions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
