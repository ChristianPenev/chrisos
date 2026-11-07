package dataloader

import (
	"context"
	"fmt"
	"net/http"
)

type contextKey[T any] struct{}

func Middleware[T any](newLoaders func() T) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextKey[T]{}, newLoaders())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func For[T any](ctx context.Context) T {
	loaders, ok := ctx.Value(contextKey[T]{}).(T)
	if !ok {
		panic(fmt.Sprintf("dataloader: no %T in context, is Middleware installed?", loaders))
	}
	return loaders
}
