package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route get /api/travel/attractions  ListAttractions
//
// ── Attractions (public/client) ──
//
// ── Attractions (public/client) ──
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: AttractionListReq
//
// Responses:
//  200: AttractionListResp

func ListAttractionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AttractionListReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewListAttractionsLogic(r.Context(), svcCtx)
		resp, err := l.ListAttractions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
