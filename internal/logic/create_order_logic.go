package logic

import (
	"context"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (*types.CreateOrderResp, error) {
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	customerID, _ := middleware.CustomerID(l.ctx)

	var travelers []*travel.Traveler
	for _, t := range req.Travelers {
		travelers = append(travelers, &travel.Traveler{
			Name: t.Name, IdType: t.IdType, IdNumber: t.IdNumber, Phone: t.Phone,
		})
	}

	resp, err := l.svcCtx.OrderClient.Create(ctx, &travel.CreateOrderRequest{
		ProductId:      req.ProductId,
		PackageId:      req.PackageId,
		Date:           req.Date,
		TimeSlot:       req.TimeSlot,
		Quantity:       int32(req.Quantity),
		CustomerEmail:  req.CustomerEmail,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		ReservationKey: req.ReservationKey,
		CustomerId:     customerID,
		Travelers:      travelers,
		Remark:         req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return toOrderResp(resp), nil
}
