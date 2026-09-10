package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCurrenciesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCurrenciesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCurrenciesLogic {
	return &ListCurrenciesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// ListCurrencies returns the global currency dictionary (public reference data).
func (l *ListCurrenciesLogic) ListCurrencies() (resp *types.CurrencyListResp, err error) {
	ctx := publicRpcContext(l.ctx)
	rpcResp, err := l.svcCtx.CurrencyClient.ListCurrencies(ctx, &travel.ListCurrenciesRequest{})
	if err != nil {
		return nil, err
	}
	items := rpcResp.GetItems()
	out := &types.CurrencyListResp{Items: make([]types.Currency, 0, len(items))}
	for _, c := range items {
		out.Items = append(out.Items, types.Currency{
			Code:           c.GetCode(),
			Symbol:         c.GetSymbol(),
			NameKey:        c.GetNameKey(),
			Decimals:       c.GetDecimals(),
			SymbolPosition: c.GetSymbolPosition(),
			IsBase:         c.GetIsBase(),
		})
	}
	return out, nil
}
