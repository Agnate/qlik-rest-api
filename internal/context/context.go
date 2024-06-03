package context

import "context"

type ContextRouteKey struct{}

func SetContextRoute(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, ContextRouteKey{}, value)
}
