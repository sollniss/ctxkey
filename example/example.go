package main

import (
	"context"
	"os"

	"github.com/sollniss/ctxkey"
	"golang.org/x/exp/slog"
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
	var logHandler slog.Handler = opts.NewJSONHandler(os.Stderr)

	// add both context keys to the log
	logHandler = someContextKey.AttachToLog("my_context_key", logHandler)
	logHandler = someOtherContextKey.AttachToLog("some_id", logHandler)

	// set the logger as default
	slog.SetDefault(slog.New(logHandler))

	// log 1
	slog.InfoCtx(ctx, "test")

	// add another value to the context, this time an int
	ctx = someOtherContextKey.WithValue(ctx, 123)

	// create a new logger with a group
	logger := slog.Default()
	logger = logger.WithGroup("test-group")

	// log 2
	logger.InfoCtx(ctx, "test2")
}
