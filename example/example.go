package main

import (
	"context"
	"os"

	"log/slog"

	"github.com/sollniss/ctxkey"
)

var someContextKey ctxkey.Key[string]
var someOtherContextKey ctxkey.Key[int]

func main() {
	ctx := context.Background()

	// add someContextKey with the value "hello world" to the context
	ctx = someContextKey.WithValue(ctx, "hello world")

	// logger options
	opts := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	// the default handler
	var logHandler slog.Handler = slog.NewJSONHandler(os.Stderr, &opts)

	// add both context keys to the log
	logHandler = someContextKey.AttachToLog("my_context_key", logHandler)
	logHandler = someOtherContextKey.AttachToLog("some_id", logHandler)

	// set the logger as default
	slog.SetDefault(slog.New(logHandler))

	// log 1
	slog.InfoContext(ctx, "test")

	// add another value to the context, this time an int
	ctx = someOtherContextKey.WithValue(ctx, 123)

	// create a new logger with a group
	logger := slog.Default()
	logger = logger.WithGroup("test-group")

	// log 2
	logger.InfoContext(ctx, "test2")

	// explicitly get value from context
	noCtxLogger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))
	some, err := someContextKey.Value(ctx)
	if err != nil {
		// key not found
		noCtxLogger.InfoContext(ctx, "explicit key", slog.String("key_not_found", err.Error()))
	} else {
		noCtxLogger.InfoContext(ctx, "explicit key", slog.String("some_key", some))
	}
}
