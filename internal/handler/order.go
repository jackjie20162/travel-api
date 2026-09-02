package handler

import (
	"errors"
	"net/http"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type OrderHandler struct { svcCtx *svc.ServiceContext }
func NewOrderHandler(svcCtx *svc.ServiceContext) *OrderHandler { return &OrderHandler{svcCtx:svcCtx} }
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderReq
	if err:=httpx.Parse(r,&req); err!=nil { httpx.Error(w,err); return }
	if req.ProductId<=0 || req.PackageId<=0 || req.Quantity<=0 || req.Date=="" || req.ReservationKey=="" { httpx.Error(w,errors.New("invalid order request")); return }
	ctx,err:=rpcContext(r); if err!=nil { httpx.Error(w,err); return }
	customerID,_:=middleware.CustomerID(ctx)
	resp,err:=h.svcCtx.OrderClient.Create(ctx,&travel.CreateOrderRequest{ProductId:req.ProductId,PackageId:req.PackageId,Date:req.Date,TimeSlot:req.TimeSlot,Quantity:int32(req.Quantity),CustomerEmail:req.CustomerEmail,ReservationKey:req.ReservationKey})
	_ = customerID
	if err!=nil { httpx.Error(w,err); return }
	httpx.OkJson(w,toOrder(resp))
}
func toOrder(o *travel.Order) CreateOrderResp { return CreateOrderResp{Id:o.Id,OrderNo:o.OrderNo,Status:o.Status,TotalAmount:o.TotalAmount,Currency:o.Currency} }
