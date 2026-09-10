// Package i18n provides a lightweight, dependency-free message bundle for the
// travel-api boundary. It translates stable backend error codes into localized
// messages driven by the request locale (Accept-Language / X-Locale).
//
// Design (see 多语言多货币改造方案 P3):
//   - Locale JSON files are embedded from ./locale and keyed by stable error code.
//   - RPC/logic errors keep their existing English messages; the API boundary maps
//     them to codes (see error.go) and localizes the response as {code, msg}.
//   - Fallbacks: requested locale -> default locale (en-US) -> empty (caller keeps raw).
package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

// Supported locales. Values match the frontend SUPPORTED_LOCALES contract.
const (
	LocaleZhCN    = "zh-CN"
	LocaleEnUS    = "en-US"
	LocaleAr      = "ar"
	DefaultLocale = LocaleEnUS
)

// localeHeader / localeHeaderAlt are the HTTP headers carrying the request locale.
const (
	localeHeader    = "Accept-Language"
	localeHeaderAlt = "X-Locale"
)

//go:embed locale/*.json
var localeFS embed.FS

var (
	bundles  = map[string]map[string]string{} // locale -> code -> message
	loadOnce sync.Once
)

// Load reads all embedded locale bundles. Safe to call multiple times.
func Load() {
	loadOnce.Do(func() {
		entries, err := localeFS.ReadDir("locale")
		if err != nil {
			logx.Errorf("i18n: failed to read embedded locale dir: %v", err)
			return
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			locale := NormalizeLocale(strings.TrimSuffix(e.Name(), ".json"))
			data, err := localeFS.ReadFile("locale/" + e.Name())
			if err != nil {
				logx.Errorf("i18n: failed to read locale file %s: %v", e.Name(), err)
				continue
			}
			var m map[string]string
			if err := json.Unmarshal(data, &m); err != nil {
				logx.Errorf("i18n: failed to parse locale file %s: %v", e.Name(), err)
				continue
			}
			bundles[locale] = m
		}
		logx.Infof("i18n: loaded %d locale bundles", len(bundles))
	})
}

// LocaleFromRequest resolves the request locale from Accept-Language (fallback X-Locale).
func LocaleFromRequest(r *http.Request) string {
	if r == nil {
		return DefaultLocale
	}
	if v := r.Header.Get(localeHeader); v != "" {
		return NormalizeLocale(v)
	}
	if v := r.Header.Get(localeHeaderAlt); v != "" {
		return NormalizeLocale(v)
	}
	return DefaultLocale
}

// NormalizeLocale maps a raw Accept-Language / X-Locale value to a supported locale.
func NormalizeLocale(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return DefaultLocale
	}
	// Take the first language tag when a list is provided (e.g. "zh-CN,zh;q=0.9").
	if i := strings.IndexAny(s, ",;"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}
	switch {
	case strings.HasPrefix(s, "zh"):
		return LocaleZhCN
	case strings.HasPrefix(s, "ar"):
		return LocaleAr
	case strings.HasPrefix(s, "en"):
		return LocaleEnUS
	default:
		return DefaultLocale
	}
}

// lookup returns the message for code in a single locale, or "" when absent.
func lookup(locale, code string) string {
	if m, ok := bundles[locale]; ok {
		if v, ok := m[code]; ok && v != "" {
			return v
		}
	}
	return ""
}

// T returns the localized message for code, falling back to the default locale.
// It returns "" when the code is unknown so callers can decide their own fallback.
func T(locale, code string) string {
	if s := lookup(locale, code); s != "" {
		return s
	}
	if locale != DefaultLocale {
		if s := lookup(DefaultLocale, code); s != "" {
			return s
		}
	}
	return ""
}

// ── context locale propagation ──

type localeContextKey struct{}

// WithLocale stores a normalized locale into the context.
func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, localeContextKey{}, locale)
}

// LocaleFrom reads the normalized locale from context, defaulting to en-US.
func LocaleFrom(ctx context.Context) string {
	if ctx == nil {
		return DefaultLocale
	}
	if v, ok := ctx.Value(localeContextKey{}).(string); ok && v != "" {
		return v
	}
	return DefaultLocale
}
