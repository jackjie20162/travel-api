package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitee.com/meinongyihe/travel-rpc/travel"
)

type ListMerchantOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMerchantOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMerchantOrdersLogic {
	return &ListMerchantOrdersLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListMerchantOrdersLogic) ListMerchantOrders(req *types.ProductListReq) (*types.OrderListResp, error) {
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	status := req.Status
	page := int32(req.Page)
	pageSize := int32(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	resp, err := l.svcCtx.OrderClient.ListMerchantOrders(ctx, &travel.MerchantOrderListRequest{
		Status:   status,
		Page:     page,
		PageSize: pageSize,
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
