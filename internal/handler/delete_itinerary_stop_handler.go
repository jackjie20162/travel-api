package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
)

// swagger:route delete /api/travel/merchant/products/{id}/itinerary-stops/{stopId}  DeleteItineraryStop
//

//

//
// Responses:
//  200: OkResp

func DeleteItineraryStopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logic.SetPathParams(r.Context(), r)
			l := logic.NewDeleteItineraryStopLogic(ctx, svcCtx)
		resp, err := l.DeleteItineraryStop()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
