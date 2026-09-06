package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReviewLogic {
	return &CreateReviewLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateReviewLogic) CreateReview(req *types.CreateReviewReq) (*types.Review, error) {
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}

	// Inject user ID for authentication
	if userID, ok := middleware.AuthUserID(l.ctx); ok && userID > 0 {
		ctx = attachCustomerID(ctx, userID)
	}

	resp, err := l.svcCtx.ReviewClient.Create(ctx, &travel.CreateReviewRequest{
		OrderNo:       req.OrderNo,
		ProductId:     req.ProductId,
		Rating:        int32(req.Rating),
		ServiceRating: int32(req.ServiceRating),
		ValueRating:   int32(req.ValueRating),
		Content:       req.Content,
		Images:        req.Images,
	})
	if err != nil {
		return nil, err
	}
	return toReviewType(resp), nil
}

func toReviewType(r *travel.Review) *types.Review {
	return &types.Review{
		Id:            r.GetId(),
		OrderNo:       r.GetOrderNo(),
		ProductId:     r.GetProductId(),
		UserId:        r.GetUserId(),
		Rating:        r.GetRating(),
		ServiceRating: r.GetServiceRating(),
		ValueRating:   r.GetValueRating(),
		Content:       r.GetContent(),
		Images:        r.GetImages(),
		ReplyContent:  r.GetReplyContent(),
		ReplyTime:     r.GetReplyTime(),
		Status:        r.GetStatus(),
		UserName:      r.GetUserName(),
		ProductName:   r.GetProductName(),
	}
}

// suppress unused import
var _ = errors.New
