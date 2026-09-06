package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (*types.Product, error) {
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	p, err := l.svcCtx.ManagementClient.CreateProduct(ctx, &travel.CreateProductRequest{
		Code: req.Code, Title: req.Title, Slug: req.Slug,
		Destination: req.Destination, Description: req.Description,
		Currency: req.Currency, MinPrice: req.MinPrice,
		Highlights: req.Highlights, CoverImage: req.CoverImage,
		Images: req.Images, VideoUrl: req.VideoUrl,
		RichContent: req.RichContent, BookingNotice: req.BookingNotice,
	})
	if err != nil {
		return nil, err
	}
	out := toProduct(p)
	return &out, nil
}
