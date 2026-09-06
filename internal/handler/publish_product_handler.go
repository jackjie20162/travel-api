package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
)

// swagger:route post /api/travel/merchant/products/{id}/publish  PublishProduct
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: PublishProductReq
//
// Responses:
//  200: Product

func PublishProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishProductReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		ctx := logic.SetPathParams(r.Context(), r)
			l := logic.NewPublishProductLogic(ctx, svcCtx)
		resp, err := l.PublishProduct(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
