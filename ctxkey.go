// copied and modified from https://github.com/alextanhongpin/go-core-microservice/blob/088138cbd567824ddec3cb39909d640c1b38a04a/http/types/contextkey.go
package ctxkey

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
)

var ErrKeyNotFound = errors.New("context key not found")

type Key[T any] struct{}

func (k Key[T]) WithValue(ctx context.Context, t T) context.Context {
	return context.WithValue(ctx, k, t)
}

func (k Key[T]) Value(ctx context.Context) (T, error) {
	t, ok := ctx.Value(k).(T)
	if !ok {
		// since the key is a struct{} we have no way to find out which key was not found
		return t, ErrKeyNotFound
	}

	return t, nil
}

func (k Key[T]) AttachToLog(key string, next slog.Handler) slog.Handler {
	return slogHandler[T]{
		logKey: key,
		ctxKey: k,
		next:   next,
	}
}

type slogHandler[T any] struct {
	logKey string
	ctxKey Key[T]
	next   slog.Handler
}

func (h slogHandler[T]) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h slogHandler[T]) Handle(ctx context.Context, record slog.Record) error {
	val, ok := ctx.Value(h.ctxKey).(T)

	if ok {
		record.AddAttrs(slog.Any(h.logKey, val))
	}

	return h.next.Handle(ctx, record)
}

func (h slogHandler[T]) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.ctxKey.AttachToLog(h.logKey, h.next.WithAttrs(attrs))
}

func (h slogHandler[T]) WithGroup(name string) slog.Handler {
	return h.ctxKey.AttachToLog(h.logKey, h.next.WithGroup(name))
}
