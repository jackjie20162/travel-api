package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductDetailLogic {
	return &ProductDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ProductDetailLogic) ProductDetail() (*types.Product, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	p, err := l.svcCtx.CatalogClient.GetProduct(ctx, &travel.ProductIdRequest{Id: id})
	if err != nil {
		return nil, err
	}
	out := toProduct(p)
	return &out, nil
}
