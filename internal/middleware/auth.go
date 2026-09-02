package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const (
	tenantIDKey  contextKey = "travel.tenant_id"
	merchantIDKey contextKey = "travel.merchant_id"
	customerIDKey contextKey = "travel.customer_id"
)

func ContextFromRequest(r *http.Request) context.Context {
	ctx := r.Context()
	ctx = context.WithValue(ctx, tenantIDKey, headerInt64(r, "X-Tenant-ID"))
	ctx = context.WithValue(ctx, merchantIDKey, headerInt64(r, "X-Merchant-ID"))
	ctx = context.WithValue(ctx, customerIDKey, headerInt64(r, "X-Customer-ID"))
	return ctx
}

func TenantID(ctx context.Context) (int64, bool) { return contextInt64(ctx, tenantIDKey) }
func MerchantID(ctx context.Context) (int64, bool) { return contextInt64(ctx, merchantIDKey) }
func CustomerID(ctx context.Context) (int64, bool) { return contextInt64(ctx, customerIDKey) }

func headerInt64(r *http.Request, name string) int64 {
	v, err := strconv.ParseInt(r.Header.Get(name), 10, 64)
	if err != nil || v <= 0 { return 0 }
	return v
}

func contextInt64(ctx context.Context, key contextKey) (int64, bool) {
	v, ok := ctx.Value(key).(int64)
	return v, ok && v > 0
}
