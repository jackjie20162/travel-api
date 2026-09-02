package handler

import (
    "errors"
    "net/http"
    "strconv"

    "gitee.com/meinongyihe/travel-api/internal/svc"
    "gitee.com/meinongyihe/travel-rpc/travel"
    "github.com/zeromicro/go-zero/rest/httpx"
)

type ManagementHandler struct { svcCtx *svc.ServiceContext }
func NewManagementHandler(svcCtx *svc.ServiceContext) *ManagementHandler { return &ManagementHandler{svcCtx:svcCtx} }

func (h *ManagementHandler) CreateProduct(w http.ResponseWriter,r *http.Request) { var req CreateProductReq; if err:=httpx.Parse(r,&req); err!=nil { httpx.Error(w,err);return }; ctx,err:=rpcContext(r);if err!=nil{httpx.Error(w,err);return}; p,err:=h.svcCtx.ManagementClient.CreateProduct(ctx,&travel.CreateProductRequest{Code:req.Code,Title:req.Title,Slug:req.Slug,Destination:req.Destination,Description:req.Description,Currency:req.Currency,MinPrice:req.MinPrice});if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,toProduct(p)) }
func (h *ManagementHandler) UpdateProduct(w http.ResponseWriter,r *http.Request) { id,err:=pathID(r,"id");if err!=nil{httpx.Error(w,err);return};var req UpdateProductReq;if err=httpx.Parse(r,&req);err!=nil{httpx.Error(w,err);return};ctx,err:=rpcContext(r);if err!=nil{httpx.Error(w,err);return};p,err:=h.svcCtx.ManagementClient.UpdateProduct(ctx,&travel.UpdateProductRequest{Id:id,Code:req.Code,Title:req.Title,Slug:req.Slug,Destination:req.Destination,Description:req.Description,Currency:req.Currency,MinPrice:req.MinPrice});if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,toProduct(p)) }
func (h *ManagementHandler) PublishProduct(w http.ResponseWriter,r *http.Request) { id,err:=pathID(r,"id");if err!=nil{httpx.Error(w,err);return};var req PublishProductReq;if err=httpx.Parse(r,&req);err!=nil{httpx.Error(w,err);return};ctx,err:=rpcContext(r);if err!=nil{httpx.Error(w,err);return};p,err:=h.svcCtx.ManagementClient.PublishProduct(ctx,&travel.PublishProductRequest{ProductId:id,Published:req.Published});if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,toProduct(p)) }
func (h *ManagementHandler) CreatePackage(w http.ResponseWriter,r *http.Request) { id,err:=pathID(r,"id");if err!=nil{httpx.Error(w,err);return};var req PackageReq;if err=httpx.Parse(r,&req);err!=nil{httpx.Error(w,err);return};ctx,err:=rpcContext(r);if err!=nil{httpx.Error(w,err);return};p,err:=h.svcCtx.ManagementClient.CreatePackage(ctx,&travel.CreatePackageRequest{ProductId:id,Code:req.Code,Name:req.Name});if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,toPackage(p)) }
func (h *ManagementHandler) ListPackages(w http.ResponseWriter,r *http.Request) { id,err:=pathID(r,"id");if err!=nil{httpx.Error(w,err);return};ctx,err:=rpcContext(r);if err!=nil{httpx.Error(w,err);return};p,err:=h.svcCtx.ManagementClient.ListPackages(ctx,&travel.PackageListRequest{ProductId:id});if err!=nil{httpx.Error(w,err);return};out:=PackageListResp{};for _,x:=range p.Items{out.Items=append(out.Items,toPackage(x))};httpx.OkJson(w,out) }
func (h *ManagementHandler) UpsertInventory(w http.ResponseWriter,r *http.Request) { id,err:=pathID(r,"id");if err!=nil{httpx.Error(w,err);return};var req InventoryUpsertReq;if err=httpx.Parse(r,&req);err!=nil{httpx.Error(w,err);return};ctx,err:=rpcContext(r);if err!=nil{httpx.Error(w,err);return};x,err:=h.svcCtx.ManagementClient.UpsertInventory(ctx,&travel.UpsertInventoryRequest{PackageId:id,Date:req.Date,TimeSlot:req.TimeSlot,Capacity:int32(req.Capacity),UnitPrice:req.UnitPrice,Currency:req.Currency,Status:req.Status});if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,toInventory(x)) }
func (h *ManagementHandler) ListInventory(w http.ResponseWriter,r *http.Request) { id,err:=pathID(r,"id");if err!=nil{httpx.Error(w,err);return};ctx,err:=rpcContext(r);if err!=nil{httpx.Error(w,err);return};x,err:=h.svcCtx.ManagementClient.ListInventory(ctx,&travel.InventoryListRequest{PackageId:id});if err!=nil{httpx.Error(w,err);return};out:=InventoryListResp{};for _,i:=range x.Items{out.Items=append(out.Items,toInventory(i))};httpx.OkJson(w,out) }

func pathID(r *http.Request,name string)(int64,error){id,err:=strconv.ParseInt(r.PathValue(name),10,64);if err!=nil||id<=0{return 0,errors.New("invalid "+name)};return id,nil}
func toPackage(p *travel.ProductPackage) Package{return Package{Id:p.Id,ProductId:p.ProductId,Code:p.Code,Name:p.Name,Status:p.Status}}
func toInventory(i *travel.InventoryItem) InventoryItem{return InventoryItem{Id:i.Id,PackageId:i.PackageId,Date:i.Date,TimeSlot:i.TimeSlot,Capacity:i.Capacity,Reserved:i.Reserved,UnitPrice:i.UnitPrice,Currency:i.Currency,Status:i.Status}}
