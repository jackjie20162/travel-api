package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductsLogic {
	return &ListProductsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListProductsLogic) ListProducts(req *types.ProductListReq) (*types.ProductListResp, error) {
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.CatalogClient.ListProducts(ctx, &travel.ProductListRequest{
		Keyword: req.Keyword, Destination: req.Destination,
		Page: int32(req.Page), PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	out := &types.ProductListResp{Total: resp.Total}
	for _, p := range resp.Items {
		out.Items = append(out.Items, toProduct(p))
	}
	// 线路没有主图时回退取行程绑定景点的图
	fillMissingProductCovers(l.ctx, l.svcCtx.DB, out.Items)
	return out, nil
}
