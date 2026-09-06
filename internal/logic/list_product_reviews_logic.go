package logic

import (
	"context"
	"errors"
	"strconv"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductReviewsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductReviewsLogic {
	return &ListProductReviewsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListProductReviewsLogic) ListProductReviews(req *types.ProductListReq) (*types.ReviewListResp, error) {
	// Extract product ID from path param
	idStr := pathVal(l.ctx, "id")
	productID, _ := strconv.ParseInt(idStr, 10, 64)
	if productID <= 0 {
		return nil, errors.New("invalid product id")
	}

	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}

	page := int32(req.Page)
	pageSize := int32(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	resp, err := l.svcCtx.ReviewClient.ListByProduct(ctx, &travel.ReviewListRequest{
		ProductId: productID,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		return nil, err
	}

	result := &types.ReviewListResp{
		Total:     resp.GetTotal(),
		AvgRating: resp.GetAvgRating(),
	}
	for _, r := range resp.GetItems() {
		result.Items = append(result.Items, types.Review{
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
		})
	}
	return result, nil
}
