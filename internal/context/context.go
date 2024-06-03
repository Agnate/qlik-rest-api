package context

import "context"

type ContextRouteKey struct{}

func GetContextRouteData(ctx context.Context) any {
	return ctx.Value(ContextRouteKey{})
}

func SetContextRouteData(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, ContextRouteKey{}, value)
}

func GetSlug(ctx context.Context, index int) string {
	data := GetContextRouteData(ctx).([]string)
	return data[index]
}
