package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gitee.com/meinongyihe/travel-api/internal/logic"
	"gitee.com/meinongyihe/travel-api/internal/svc"
)

// swagger:route get /api/travel/currencies  ListCurrencies
//
// ── Currency & Exchange Rate (public reference data) ──
//
// ── Currency & Exchange Rate (public reference data) ──
//
// Responses:
//  200: CurrencyListResp

func ListCurrenciesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewListCurrenciesLogic(r.Context(), svcCtx)
		resp, err := l.ListCurrencies()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
