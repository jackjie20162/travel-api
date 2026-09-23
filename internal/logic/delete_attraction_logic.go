package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/attraction"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAttractionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAttractionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAttractionLogic {
	return &DeleteAttractionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *DeleteAttractionLogic) DeleteAttraction() (resp *types.OkResp, err error) {
	id, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	merchantID, _ := middleware.MerchantID(l.ctx)

	_, err = l.svcCtx.DB.Attraction.Delete().
		Where(attraction.IDEQ(int(id)), attraction.MerchantIDEQ(merchantID)).
		Exec(l.ctx)
	if err != nil {
		return nil, err
	}
	return &types.OkResp{Ok: true, Msg: "deleted"}, nil
}
