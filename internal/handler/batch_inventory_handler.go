package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route get /api/travel/inventory/packages/{packageId}/batch  BatchInventory
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: InventoryBatchReq
//
// Responses:
//  200: InventoryListResp

func BatchInventoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InventoryBatchReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		ctx := logic.SetPathParams(r.Context(), r)
			l := logic.NewBatchInventoryLogic(ctx, svcCtx)
		resp, err := l.BatchInventory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
