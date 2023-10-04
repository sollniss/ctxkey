# ctxkey

Provides generic context keys with getter and setter methods. Also provides a function to automatically log the context key.

# Example


```go
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
```

prints

```
{"time":"2023-10-04T11:37:00.6478478+09:00","level":"INFO","msg":"test","my_context_key":"hello world"}
{"time":"2023-10-04T11:37:00.6737874+09:00","level":"INFO","msg":"test2","test-group":{"some_id":123,"my_context_key":"hello world"}}
{"time":"2023-10-04T11:37:00.6737874+09:00","level":"INFO","msg":"explicit key","some_key":"hello world"}
```

Note that the context keys will be attached to the log last. So they are inside the group we defined.