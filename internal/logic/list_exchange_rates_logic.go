package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListExchangeRatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListExchangeRatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListExchangeRatesLogic {
	return &ListExchangeRatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// ListExchangeRates returns the latest exchange rates for a base currency (default AED).
func (l *ListExchangeRatesLogic) ListExchangeRates(req *types.ExchangeRateListReq) (resp *types.ExchangeRateListResp, err error) {
	ctx := publicRpcContext(l.ctx)
	rpcResp, err := l.svcCtx.CurrencyClient.ListExchangeRates(ctx, &travel.ListExchangeRatesRequest{Base: req.Base})
	if err != nil {
		return nil, err
	}
	items := rpcResp.GetItems()
	out := &types.ExchangeRateListResp{Items: make([]types.ExchangeRate, 0, len(items))}
	for _, r := range items {
		out.Items = append(out.Items, types.ExchangeRate{
			BaseCurrency:   r.GetBaseCurrency(),
			TargetCurrency: r.GetTargetCurrency(),
			Rate:           float64(r.GetRateMicro()) / 1e6,
			RateMicro:      r.GetRateMicro(),
			EffectiveAt:    r.GetEffectiveAt(),
		})
	}
	return out, nil
}
