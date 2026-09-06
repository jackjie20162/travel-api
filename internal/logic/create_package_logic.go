package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePackageLogic {
	return &CreatePackageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreatePackageLogic) CreatePackage(req *types.PackageReq) (*types.Package, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	p, err := l.svcCtx.ManagementClient.CreatePackage(ctx, &travel.CreatePackageRequest{
		ProductId: id, Code: req.Code, Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	out := toPackage(p)
	return &out, nil
}
