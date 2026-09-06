package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishProductLogic {
	return &PublishProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PublishProductLogic) PublishProduct(req *types.PublishProductReq) (*types.Product, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	p, err := l.svcCtx.ManagementClient.PublishProduct(ctx, &travel.PublishProductRequest{
		ProductId: id, Published: req.Published,
	})
	if err != nil {
		return nil, err
	}
	out := toProduct(p)
	return &out, nil
}
