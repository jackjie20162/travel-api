package logic

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/pathvar"
	"google.golang.org/grpc/metadata"
)

// ── gRPC context builders ──

func rpcContext(ctx context.Context) (context.Context, error) {
	tenant, ok := middleware.TenantID(ctx)
	if !ok {
		return nil, errors.New("missing tenant context")
	}
	merchant, ok := middleware.MerchantID(ctx)
	if !ok {
		return nil, errors.New("missing merchant context")
	}
	pairs := []string{
		"x-tenant-id", strconv.FormatInt(tenant, 10),
		"x-merchant-id", strconv.FormatInt(merchant, 10),
	}
	// Pass user_id (归属人) if available from authenticated user
	if userID, ok := middleware.AuthUserID(ctx); ok && userID > 0 {
		pairs = append(pairs, "x-customer-id", strconv.FormatInt(userID, 10))
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(pairs...)), nil
}

func customerRpcContext(ctx context.Context, customerID int64) (context.Context, error) {
	tenant, ok := middleware.TenantID(ctx)
	if !ok {
		return nil, errors.New("missing tenant context")
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"x-tenant-id", strconv.FormatInt(tenant, 10),
		"x-customer-id", strconv.FormatInt(customerID, 10),
	)), nil
}

func userRpcContext(ctx context.Context, userID int64) (context.Context, error) {
	tenant, ok := middleware.TenantID(ctx)
	if !ok {
		return nil, errors.New("missing tenant context")
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"x-tenant-id", strconv.FormatInt(tenant, 10),
		"x-user-id", strconv.FormatInt(userID, 10),
	)), nil
}

// attachCustomerID adds x-customer-id to existing gRPC metadata context.
func attachCustomerID(ctx context.Context, userID int64) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	md.Set("x-customer-id", strconv.FormatInt(userID, 10))
	return metadata.NewOutgoingContext(ctx, md)
}

// ── path param helpers via context ──

func pathVal(ctx context.Context, name string) string {
	v, _ := ctx.Value(pathKey(name)).(string)
	return v
}

func pathInt64(ctx context.Context, name string) (int64, error) {
	s := pathVal(ctx, name)
	if s == "" {
		return 0, errors.New("missing path param: " + name)
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v <= 0 {
		return 0, errors.New("invalid " + name)
	}
	return v, nil
}

type pathParamKey string

func pathKey(name string) pathParamKey { return pathParamKey(name) }

// SetPathParams extracts path variables from the request and stores them in context.
func SetPathParams(ctx context.Context, r *http.Request) context.Context {
	// go-zero stores path params via pathvar.Vars(r)
	for k, v := range pathVars(r) {
		ctx = context.WithValue(ctx, pathKey(k), v)
	}
	return ctx
}

// pathVars wraps go-zero's pathvar.Vars to extract route parameters.
func pathVars(r *http.Request) map[string]string {
	return pathvar.Vars(r)
}

// ── proto → types converters ──

func toProduct(p *travel.Product) types.Product {
	return types.Product{
		Id: p.Id, TenantId: p.TenantId, MerchantId: p.MerchantId,
		Code: p.Code, Title: p.Title, Slug: p.Slug,
		Destination: p.Destination, Description: p.Description,
		Currency: p.Currency, MinPrice: p.MinPrice, Status: p.Status,
		Highlights: p.Highlights, CoverImage: p.CoverImage,
		Images: p.Images, VideoUrl: p.VideoUrl,
		RichContent: p.RichContent, BookingNotice: p.BookingNotice,
	}
}

func toPackage(p *travel.ProductPackage) types.Package {
	return types.Package{
		Id: p.Id, ProductId: p.ProductId, Code: p.Code,
		Name: p.Name, Status: p.Status,
		PricingMode: p.PricingMode, InventoryMode: p.InventoryMode,
		MinOrderQty: int(p.MinOrderQty), SellCurrency: p.SellCurrency,
		CostCurrency: p.CostCurrency, GroupPrices: p.GroupPrices,
		TierPrices: p.TierPrices,
	}
}

func toInventoryItem(i *travel.InventoryItem) types.InventoryItem {
	return types.InventoryItem{
		Id: i.Id, PackageId: i.PackageId, Date: i.Date,
		TimeSlot: i.TimeSlot, Capacity: i.Capacity, Reserved: i.Reserved,
		UnitPrice: i.UnitPrice, Currency: i.Currency, Status: i.Status,
		InventoryMode: i.InventoryMode, TotalCapacity: i.TotalCapacity,
		IsOpen: i.IsOpen,
	}
}

func toOrderResp(o *travel.Order) *types.CreateOrderResp {
	resp := &types.CreateOrderResp{
		Id: o.Id, OrderNo: o.OrderNo, Status: o.Status,
		TotalAmount: o.TotalAmount, Currency: o.Currency,
		CustomerEmail: o.GetCustomerEmail(), PaymentStatus: o.GetPaymentStatus(),
		CustomerId: o.GetCustomerId(), CustomerName: o.GetCustomerName(),
		CustomerPhone: o.GetCustomerPhone(), Remark: o.GetRemark(),
		ProductName: o.GetProductName(), PackageName: o.GetPackageName(),
		ServiceDate: o.GetServiceDate(), TimeSlot: o.GetTimeSlot(),
		RejectReason: o.GetRejectReason(), VerifiedAt: o.GetVerifiedAt(),
	}
	for _, t := range o.GetTravelers() {
		resp.Travelers = append(resp.Travelers, types.Traveler{
			Id: t.GetId(), Name: t.GetName(), IdType: t.GetIdType(),
			IdNumber: t.GetIdNumber(), Phone: t.GetPhone(),
		})
	}
	for _, item := range o.GetItems() {
		resp.Items = append(resp.Items, types.OrderItem{
			Id: item.GetId(), ProductId: item.GetProductId(), PackageId: item.GetPackageId(),
			TravelerId: item.GetTravelerId(), Quantity: item.GetQuantity(),
			UnitPrice: item.GetUnitPrice(), TotalAmount: item.GetTotalAmount(),
			ServiceDate: item.GetServiceDate(), TimeSlot: item.GetTimeSlot(),
			ProductCode: item.GetProductCode(), ProductName: item.GetProductName(),
			PackageCode: item.GetPackageCode(), PackageName: item.GetPackageName(),
		})
	}
	return resp
}

func toTravelers(ts []*travel.Traveler) []types.Traveler {
	if len(ts) == 0 {
		return nil
	}
	out := make([]types.Traveler, len(ts))
	for i, t := range ts {
		out[i] = types.Traveler{
			Id: t.GetId(), Name: t.GetName(), IdType: t.GetIdType(),
			IdNumber: t.GetIdNumber(), Phone: t.GetPhone(),
		}
	}
	return out
}

func toPaymentResp(p *travel.Payment) *types.PaymentResp {
	return &types.PaymentResp{
		Id: p.Id, PaymentNo: p.PaymentNo, OrderNo: p.OrderNo,
		Provider: p.Provider, ProviderPaymentId: p.ProviderPaymentId,
		Amount: p.Amount, Currency: p.Currency, Status: p.Status,
	}
}

func toUserProfile(u *travel.User) types.UserProfile {
	return types.UserProfile{
		Id: u.Id, Username: u.Username, Nickname: u.Nickname,
		Email: u.Email, Mobile: u.Mobile, Avatar: u.Avatar,
		Level: u.Level, Gender: u.Gender, Birthday: u.Birthday, Bio: u.Bio,
		Money: u.Money, Score: u.Score, Jointime: u.Jointime,
		Logintime: u.Logintime, Loginip: u.Loginip,
		Status: u.Status, Createtime: u.Createtime,
	}
}

func toItineraryStop(s *travel.ItineraryStop) types.ItineraryStop {
	return types.ItineraryStop{
		Id: s.Id, ProductId: s.ProductId, StopType: s.StopType,
		Title: s.Title, Description: s.Description, Sequence: int(s.Sequence),
		LocationName: s.LocationName, Address: s.Address,
		Latitude: s.Latitude, Longitude: s.Longitude,
		DurationMinutes: int(s.DurationMinutes), TransportType: s.TransportType,
		StartTime: s.StartTime, EndTime: s.EndTime,
		IsPickup: s.IsPickup, IsDropoff: s.IsDropoff,
		AgreementNoShopping: s.AgreementNoShopping, AgreementAdjustable: s.AgreementAdjustable,
		PoiId: s.PoiId, PoiName: s.PoiName, IsEntering: s.IsEntering,
		DurationMode: s.DurationMode, DurationHours: int(s.DurationHours),
		ActivityFeatures: s.ActivityFeatures,
		PickupLocation: s.PickupLocation, PickupAddress: s.PickupAddress,
		PickupLatitude: s.PickupLatitude, PickupLongitude: s.PickupLongitude,
		DropoffLocation: s.DropoffLocation, DropoffAddress: s.DropoffAddress,
		DropoffLatitude: s.DropoffLatitude, DropoffLongitude: s.DropoffLongitude,
	}
}

// suppress unused import warnings
var _ = logx.Info
var _ = svc.ServiceContext{}
