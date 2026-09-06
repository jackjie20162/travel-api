package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplyReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyReviewLogic {
	return &ReplyReviewLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ReplyReviewLogic) ReplyReview(req *types.ReplyReviewReq) (*types.Review, error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}

	ctx, rpcErr := rpcContext(l.ctx)
	if rpcErr != nil {
		return nil, rpcErr
	}

	resp, rpcErr := l.svcCtx.ReviewClient.Reply(ctx, &travel.ReplyReviewRequest{
		Id:           id,
		ReplyContent: req.ReplyContent,
	})
	if rpcErr != nil {
		return nil, rpcErr
	}
	return toReviewType(resp), nil
}
