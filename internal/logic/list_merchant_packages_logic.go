package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListMerchantPackagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMerchantPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMerchantPackagesLogic {
	return &ListMerchantPackagesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListMerchantPackagesLogic) ListMerchantPackages() (*types.PackageListResp, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.ManagementClient.ListPackages(ctx, &travel.PackageListRequest{ProductId: id})
	if err != nil {
		return nil, err
	}
	out := &types.PackageListResp{}
	for _, p := range resp.Items {
		out.Items = append(out.Items, toPackage(p))
	}
	return out, nil
}
