package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
)

// swagger:route post /api/travel/merchant/products/{id}/sync-shopify  SyncProductToShopify
//
// ── Shopify Sales Channel ──
//
// ── Shopify Sales Channel ──
//
// Responses:
//  200: OkResp

func SyncProductToShopifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewSyncProductToShopifyLogic(r.Context(), svcCtx)
		resp, err := l.SyncProductToShopify()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
