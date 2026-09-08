package logic

import (
	"context"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateItineraryStopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateItineraryStopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateItineraryStopLogic {
	return &UpdateItineraryStopLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateItineraryStopLogic) UpdateItineraryStop(req *types.ItineraryStopReq) (*types.ItineraryStop, error) {
	productId, err := pathInt64(l.ctx, "id")
	if err != nil {
		return nil, err
	}
	stopId, err := pathInt64(l.ctx, "stopId")
	if err != nil {
		return nil, err
	}
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	s, err := l.svcCtx.ManagementClient.UpdateItineraryStop(ctx, &travel.UpdateItineraryStopRequest{
		Id: stopId, ProductId: productId, StopType: req.StopType, Title: req.Title,
		Description: req.Description, Sequence: int32(req.Sequence),
		LocationName: req.LocationName, Address: req.Address,
		Latitude: req.Latitude, Longitude: req.Longitude,
		DurationMinutes: int32(req.DurationMinutes), TransportType: req.TransportType,
		StartTime: req.StartTime, EndTime: req.EndTime,
		IsPickup: req.IsPickup, IsDropoff: req.IsDropoff,
		AgreementNoShopping: req.AgreementNoShopping, AgreementAdjustable: req.AgreementAdjustable,
		PoiId: req.PoiId, PoiName: req.PoiName, IsEntering: req.IsEntering,
		DurationMode: req.DurationMode, DurationHours: int32(req.DurationHours),
		ActivityFeatures: req.ActivityFeatures,
		PickupLocation: req.PickupLocation, PickupAddress: req.PickupAddress,
		PickupLatitude: req.PickupLatitude, PickupLongitude: req.PickupLongitude,
		DropoffLocation: req.DropoffLocation, DropoffAddress: req.DropoffAddress,
		DropoffLatitude: req.DropoffLatitude, DropoffLongitude: req.DropoffLongitude,
		TypeParams: req.TypeParams,
	})
	if err != nil {
		return nil, err
	}
	out := toItineraryStop(s)
	return &out, nil
}
