package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderReviewLogic {
	return &GetOrderReviewLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetOrderReviewLogic) GetOrderReview() (*types.Review, error) {
	orderNo := pathVal(l.ctx, "orderNo")
	if orderNo == "" {
		return nil, errors.New("invalid order number")
	}

	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.ReviewClient.GetByOrder(ctx, &travel.OrderNoReviewRequest{OrderNo: orderNo})
	if err != nil {
		return nil, err
	}
	return toReviewType(resp), nil
}
