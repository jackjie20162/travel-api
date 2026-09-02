package handler

import (
	"net/http"
	"strconv"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/metadata"
)

type ProductHandler struct { svcCtx *svc.ServiceContext }

func NewProductHandler(svcCtx *svc.ServiceContext) *ProductHandler { return &ProductHandler{svcCtx: svcCtx} }

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	var req ProductListReq
	if err := httpx.Parse(r, &req); err != nil { httpx.Error(w, err); return }
	ctx, err := rpcContext(r)
	if err != nil { httpx.Error(w, err); return }
	resp, err := h.svcCtx.CatalogClient.ListProducts(ctx, &travel.ProductListRequest{Keyword:req.Keyword, Destination:req.Destination, Page:int32(req.Page), PageSize:int32(req.PageSize)})
	if err != nil { httpx.Error(w, err); return }
	items := make([]Product, 0, len(resp.Items))
	for _, p := range resp.Items { items = append(items, toProduct(p)) }
	httpx.OkJson(w, ProductListResp{Items:items, Total:resp.Total})
}

func (h *ProductHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 { httpx.Error(w, err); return }
	ctx, err := rpcContext(r)
	if err != nil { httpx.Error(w, err); return }
	p, err := h.svcCtx.CatalogClient.GetProduct(ctx, &travel.ProductIdRequest{Id:id})
	if err != nil { httpx.Error(w, err); return }
	httpx.OkJson(w, toProduct(p))
}

func rpcContext(r *http.Request) (context.Context, error) {
	ctx := middleware.ContextFromRequest(r)
	tenant, ok := middleware.TenantID(ctx); if !ok { return nil, errors.New("missing tenant context") }
	merchant, ok := middleware.MerchantID(ctx); if !ok { return nil, errors.New("missing merchant context") }
	return metadata.AppendToOutgoingContext(ctx, "x-tenant-id", strconv.FormatInt(tenant,10), "x-merchant-id", strconv.FormatInt(merchant,10)), nil
}

func toProduct(p *travel.Product) Product { return Product{Id:p.Id,TenantId:p.TenantId,MerchantId:p.MerchantId,Code:p.Code,Title:p.Title,Slug:p.Slug,Destination:p.Destination,Description:p.Description,Currency:p.Currency,MinPrice:p.MinPrice,Status:p.Status} }
