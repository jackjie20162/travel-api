package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

const (
	tenantIDKey   contextKey = "travel.tenant_id"
	merchantIDKey contextKey = "travel.merchant_id"
	customerIDKey contextKey = "travel.customer_id"
	authUserIDKey contextKey = "travel.auth_user_id"
)

func ContextFromRequest(r *http.Request) context.Context {
	ctx := r.Context()
	ctx = context.WithValue(ctx, tenantIDKey, headerInt64(r, "X-Tenant-ID"))
	ctx = context.WithValue(ctx, merchantIDKey, headerInt64(r, "X-Merchant-ID"))
	ctx = context.WithValue(ctx, customerIDKey, headerInt64(r, "X-Customer-ID"))
	return ctx
}

// WithAuthUserID stores the authenticated user ID into context.
func WithAuthUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, authUserIDKey, userID)
}

// AuthUserID retrieves the authenticated user ID from context.
func AuthUserID(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(authUserIDKey).(int64)
	return v, ok && v > 0
}

func TenantID(ctx context.Context) (int64, bool)   { return contextInt64(ctx, tenantIDKey) }
func MerchantID(ctx context.Context) (int64, bool) { return contextInt64(ctx, merchantIDKey) }
func CustomerID(ctx context.Context) (int64, bool) { return contextInt64(ctx, customerIDKey) }

// ExtractToken extracts the Bearer token from the Authorization header.
func ExtractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	// Fallback: also accept the token from X-Token header
	if t := r.Header.Get("X-Token"); t != "" {
		return t
	}
	return ""
}

func headerInt64(r *http.Request, name string) int64 {
	v, err := strconv.ParseInt(r.Header.Get(name), 10, 64)
	if err != nil || v <= 0 {
		return 0
	}
	return v
}

func contextInt64(ctx context.Context, key contextKey) (int64, bool) {
	v, ok := ctx.Value(key).(int64)
	return v, ok && v > 0
}
