package config

import "context"

type contextKey int

const configContextKey contextKey = iota + 1

func ContextWithConfigOptions(ctx context.Context, opts ...Option) context.Context {
	return context.WithValue(ctx, configContextKey, opts)
}

func ConfigOptionsFromContext(ctx context.Context) []Option {
	opts, ok := ctx.Value(configContextKey).([]Option)
	if !ok {
		return []Option{}
	}
	return opts
}
