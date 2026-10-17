package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
)

// swagger:route get /api/travel/attractions/{id}/products  ListAttractionProducts
//
// ── 公开景点关联线路列表（行程中包含该景点的已发布产品） ──
//
// Responses:
//  200: ProductListResp

func ListAttractionProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logic.SetPathParams(r.Context(), r)
		l := logic.NewListAttractionProductsLogic(ctx, svcCtx)
		resp, err := l.ListAttractionProducts()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
