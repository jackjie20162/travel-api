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

type ListMyOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyOrdersLogic {
	return &ListMyOrdersLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListMyOrdersLogic) ListMyOrders(req *types.OrderListReq) (*types.OrderListResp, error) {
	userID, ok := middleware.AuthUserID(l.ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	ctx, err := userRpcContext(l.ctx, userID)
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
	resp, err := l.svcCtx.OrderClient.ListCustomerOrders(ctx, &travel.CustomerOrderListRequest{
		CustomerId: userID,
		Status:     req.Status,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		return nil, err
	}
	out := &types.OrderListResp{Total: resp.GetTotal()}
	for _, o := range resp.GetItems() {
		out.Items = append(out.Items, *toOrderResp(o))
	}
	return out, nil
}
