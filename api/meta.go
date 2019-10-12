package api

import "context"

// A ContextID represents context key for middleware.
type ContextID struct {
	name string
}

// ContextKeyName is the key name for storing information in request context
var ContextKeyName = &ContextID{name: "sejastip"}

//MetaFromContext is function to get metainfo from context
func MetaFromContext(ctx context.Context) ResourceClaims {
	return ctx.Value(ContextKeyName).(ResourceClaims)
}
