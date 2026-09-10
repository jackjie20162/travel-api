package i18n

import (
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// messageToCode maps known backend (RPC/logic) English messages to stable error
// codes. Messages not listed here fall back to a generic code derived from the
// gRPC status code (see grpcCodeToGeneric), so every error still gets localized.
var messageToCode = map[string]string{
	// ── NotFound ──
	"product not found":        "PRODUCT_NOT_FOUND",
	"package not found":        "PACKAGE_NOT_FOUND",
	"order not found":          "ORDER_NOT_FOUND",
	"inventory not found":      "INVENTORY_NOT_FOUND",
	"review not found":         "REVIEW_NOT_FOUND",
	"user not found":           "USER_NOT_FOUND",
	"itinerary stop not found": "ITINERARY_STOP_NOT_FOUND",
	// ── ResourceExhausted ──
	"insufficient inventory": "INVENTORY_INSUFFICIENT",
	// ── AlreadyExists ──
	"username already taken":               "USERNAME_ALREADY_TAKEN",
	"mobile already registered":            "MOBILE_ALREADY_REGISTERED",
	"this order has already been reviewed": "ORDER_ALREADY_REVIEWED",
	// ── Unauthenticated ──
	"invalid password":     "INVALID_PASSWORD",
	"invalid old password": "INVALID_OLD_PASSWORD",
	// ── PermissionDenied ──
	"account is disabled": "ACCOUNT_DISABLED",
	// ── FailedPrecondition (order/product state) ──
	"order is not in pending verify status":        "ORDER_NOT_PENDING_VERIFY",
	"order is not in pending acceptance status":    "ORDER_NOT_PENDING_ACCEPTANCE",
	"order is not in pending refund status":        "ORDER_NOT_PENDING_REFUND",
	"order status does not allow refund request":   "ORDER_STATUS_NO_REFUND",
	"capacity cannot be below reserved quantity":   "CAPACITY_BELOW_RESERVED",
	"product requires at least one active package": "PRODUCT_NEEDS_ACTIVE_PACKAGE",
	"published product cannot add packages":        "PUBLISHED_PRODUCT_NO_PACKAGE",
	// ── Internal ──
	"failed to fetch updated order": "FAILED_FETCH_UPDATED_ORDER",
	// ── travel-api logic (plain errors.New) ──
	"missing tenant context":   "MISSING_TENANT_CONTEXT",
	"missing merchant context": "MISSING_MERCHANT_CONTEXT",
	"unauthorized":             "UNAUTHORIZED",
	"not yet supported":        "NOT_SUPPORTED",
	// ── travel-api validation / payment (plain errors.New) ──
	"invalid order number":                              "INVALID_ORDER_NUMBER",
	"invalid payment number":                            "INVALID_PAYMENT_NUMBER",
	"invalid product id":                                "INVALID_PRODUCT_ID",
	"invalid packageId":                                 "INVALID_PACKAGE_ID",
	"only paypal is enabled":                            "ONLY_PAYPAL_ENABLED",
	"paypal client secret is not configured":            "PAYPAL_NOT_CONFIGURED",
	"missing state parameter":                           "MISSING_STATE_PARAM",
	"missing user context":                              "MISSING_USER_CONTEXT",
	"orderNo, provider and idempotencyKey are required": "PAYMENT_PARAMS_REQUIRED",
	// ── travel-api merchant/config/inventory validation (plain errors.New) ──
	"missing packageId":                   "MISSING_PACKAGE_ID",
	"tenant context is required":          "MISSING_TENANT_CONTEXT",
	"merchant context is required":        "MISSING_MERCHANT_CONTEXT",
	"missing configType":                  "MISSING_CONFIG_TYPE",
	"configType is required":              "MISSING_CONFIG_TYPE",
	"items is required":                   "CONFIG_ITEMS_REQUIRED",
	"configKey is required for all items": "CONFIG_KEY_REQUIRED",
	// ── verification-code flow (register/login) ──
	"图形验证码错误或已过期": "CAPTCHA_INVALID",
	"邮件验证码错误或已过期": "EMAIL_CODE_INVALID",
}

// grpcCodeToGeneric maps a gRPC status code to a generic error code used when the
// specific message is unknown or dynamic (e.g. err.Error() passthroughs).
var grpcCodeToGeneric = map[codes.Code]string{
	codes.Canceled:           "ABORTED",
	codes.Unknown:            "UNKNOWN",
	codes.InvalidArgument:    "INVALID_ARGUMENT",
	codes.DeadlineExceeded:   "UNAVAILABLE",
	codes.NotFound:           "NOT_FOUND",
	codes.AlreadyExists:      "ALREADY_EXISTS",
	codes.PermissionDenied:   "PERMISSION_DENIED",
	codes.ResourceExhausted:  "RESOURCE_EXHAUSTED",
	codes.FailedPrecondition: "FAILED_PRECONDITION",
	codes.Aborted:            "ABORTED",
	codes.OutOfRange:         "INVALID_ARGUMENT",
	codes.Unimplemented:      "NOT_SUPPORTED",
	codes.Internal:           "INTERNAL_ERROR",
	codes.Unavailable:        "UNAVAILABLE",
	codes.DataLoss:           "INTERNAL_ERROR",
	codes.Unauthenticated:    "UNAUTHORIZED",
}

// errorBody is the JSON response shape {code, msg} consumed by the frontend
// request() wrapper (it throws Error(msg) when code !== 0).
type errorBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// TranslateError resolves err into a numeric code and a localized message for the
// given locale. Resolution order for gRPC errors: message-as-code / mapped code ->
// generic gRPC-code message -> original message. The numeric code preserves the
// gRPC code for continuity with the previous default handler.
func TranslateError(locale string, err error) (int, string) {
	if err == nil {
		return 0, ""
	}

	if st, ok := status.FromError(err); ok {
		numCode := int(st.Code())
		msg := st.Message()
		if s := resolveMessage(locale, msg); s != "" {
			return numCode, s
		}
		if generic, found := grpcCodeToGeneric[st.Code()]; found {
			if s := T(locale, generic); s != "" {
				return numCode, s
			}
		}
		return numCode, msg
	}

	// Plain (non-gRPC) error: logic errors.New, validation, etc.
	msg := err.Error()
	if s := resolveMessage(locale, msg); s != "" {
		return http.StatusBadRequest, s
	}
	return http.StatusBadRequest, msg
}

// resolveMessage tries the raw message as a stable code first (new code returns
// e.g. errors.New("CAPTCHA_INVALID")), then via the legacy message->code map.
// It returns "" when no localization is available so callers can fall back.
func resolveMessage(locale, msg string) string {
	if msg == "" {
		return ""
	}
	if s := T(locale, msg); s != "" {
		return s
	}
	if code, found := messageToCode[msg]; found {
		if s := T(locale, code); s != "" {
			return s
		}
	}
	return ""
}

// ErrorHandler is registered via httpx.SetErrorHandlerCtx. It localizes any error
// and always responds HTTP 200 with a {code, msg} body so the frontend surfaces the
// localized message instead of a raw HTTP error string.
func ErrorHandler(ctx context.Context, err error) (int, any) {
	numCode, msg := TranslateError(LocaleFrom(ctx), err)
	if numCode == 0 {
		numCode = http.StatusInternalServerError
	}
	return http.StatusOK, errorBody{Code: numCode, Msg: msg}
}
