# ctxkey

Provides generic context keys with getter and setter methods. Also provides a function to automatically log the context key.

# Example


```go
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

	// set the log handler as default
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
```

prints

```
{"time":"2023-04-12T02:19:57.9305901+09:00","level":"INFO","msg":"test","my_context_key":"hello world"}
{"time":"2023-04-12T02:19:57.9500553+09:00","level":"INFO","msg":"test2","test-group":{"some_id":123,"my_context_key":"hello world"}}
```

Note that the context keys will be attached to the log last. So they are inside the group we defined.