package handler

import (
	"net/http"
	"errors"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type InventoryHandler struct { svcCtx *svc.ServiceContext }
func NewInventoryHandler(svcCtx *svc.ServiceContext) *InventoryHandler { return &InventoryHandler{svcCtx:svcCtx} }
func (h *InventoryHandler) Check(w http.ResponseWriter, r *http.Request) {
	var req InventoryCheckReq
	if err:=httpx.Parse(r,&req); err!=nil { httpx.Error(w,err); return }
	if req.PackageId<=0 || req.Quantity<=0 || req.Date=="" { httpx.Error(w,errors.New("invalid inventory request")); return }
	ctx,err:=rpcContext(r); if err!=nil { httpx.Error(w,err); return }
	resp,err:=h.svcCtx.InventoryClient.Check(ctx,&travel.InventoryRequest{PackageId:req.PackageId,Date:req.Date,TimeSlot:req.TimeSlot,Quantity:int32(req.Quantity)})
	if err!=nil { httpx.Error(w,err); return }
	httpx.OkJson(w,InventoryCheckResp{Available:resp.Available,Remaining:resp.Remaining,UnitPrice:resp.UnitPrice,Currency:resp.Currency})
}
