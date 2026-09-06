package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductReq) (*types.Product, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	p, err := l.svcCtx.ManagementClient.UpdateProduct(ctx, &travel.UpdateProductRequest{
		Id: id, Code: req.Code, Title: req.Title, Slug: req.Slug,
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
