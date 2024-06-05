package context

import (
	"context"
	"errors"
)

type ContextRouteKey struct{}

func GetContextRouteData(ctx context.Context) any {
	return ctx.Value(ContextRouteKey{})
}

func SetContextRouteData(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, ContextRouteKey{}, value)
}

// Get the slug from the route URL based on index.
func GetSlug(ctx context.Context, index int) (string, error) {
	data := GetContextRouteData(ctx).([]string)
	// data, ok to ensure we can convert to []string
	if len(data) <= index {
		return "", errors.New("invalid slug")
	}
	return data[index], nil
}
