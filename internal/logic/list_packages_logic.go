package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListPackagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPackagesLogic {
	return &ListPackagesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListPackagesLogic) ListPackages() (*types.PackageListResp, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.CatalogClient.ListPackages(ctx, &travel.PackageListRequest{ProductId: id})
	if err != nil {
		return nil, err
	}
	out := &types.PackageListResp{}
	for _, p := range resp.Items {
		out.Items = append(out.Items, toPackage(p))
	}
	return out, nil
}
