package api

import (
	"context"
	"net/http"
	"strings"

	"sejastip.id/api/entity"
)

// A ContextID represents context key for middleware.
type ContextID struct {
	name string
}

// ContextKeyName is the key name for storing information in request context
var ContextKeyName = &ContextID{name: "sejastip"}

//MetaFromContext is function to get metainfo from context
func MetaFromContext(ctx context.Context) entity.ResourceClaims {
	return ctx.Value(ContextKeyName).(entity.ResourceClaims)
}

func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

func GetPlatform(userAgent string) string {
	lowerUserAgent := strings.ToLower(userAgent)
	switch {
	case strings.Contains(lowerUserAgent, "android"):
		return "android"
	case strings.Contains(lowerUserAgent, "ios"):
		return "ios"
	default:
		return "other"
	}
}
