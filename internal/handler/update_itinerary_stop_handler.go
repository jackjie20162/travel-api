package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route put /api/travel/merchant/products/{id}/itinerary-stops/{stopId}  UpdateItineraryStop
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ItineraryStopReq
//
// Responses:
//  200: ItineraryStop

func UpdateItineraryStopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ItineraryStopReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		ctx := logic.SetPathParams(r.Context(), r)
			l := logic.NewUpdateItineraryStopLogic(ctx, svcCtx)
		resp, err := l.UpdateItineraryStop(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
