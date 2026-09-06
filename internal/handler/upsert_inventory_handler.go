package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route post /api/travel/merchant/packages/{id}/inventory  UpsertInventory
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: InventoryUpsertReq
//
// Responses:
//  200: InventoryItem

func UpsertInventoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InventoryUpsertReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		ctx := logic.SetPathParams(r.Context(), r)
			l := logic.NewUpsertInventoryLogic(ctx, svcCtx)
		resp, err := l.UpsertInventory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
