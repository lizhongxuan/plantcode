package log

import (
	"ai-dev-platform/internal/requestid"
	"context"
	"fmt"
)

func DebugId(ctx context.Context, args ...interface{}) { // nolint
	logger.Debug(requestid.String(ctx) + " " + fmt.Sprintln(args...))
}

func DebugfId(ctx context.Context, format string, args ...interface{}) { // nolint
	logger.Debug(requestid.String(ctx) + " " + fmt.Sprintf(format, args...))
}

func InfoId(ctx context.Context, args ...interface{}) { // nolint
	logger.Info(requestid.String(ctx) + " " + fmt.Sprintln(args...))
}

func InfofId(ctx context.Context, format string, args ...interface{}) { // nolint
	logger.Info(requestid.String(ctx) + " " + fmt.Sprintf(format, args...))
}

func WarnId(ctx context.Context, args ...interface{}) { // nolint
	logger.Warn(requestid.String(ctx) + " " + fmt.Sprintln(args...))
}

func WarnfId(ctx context.Context, format string, args ...interface{}) { // nolint
	logger.Warn(requestid.String(ctx) + " " + fmt.Sprintf(format, args...))
}

func ErrorId(ctx context.Context, args ...interface{}) { // nolint
	logger.Error(requestid.String(ctx) + " " + fmt.Sprintln(args...))
}

func ErrorfId(ctx context.Context, format string, args ...interface{}) { // nolint
	logger.Error(requestid.String(ctx) + " " + fmt.Sprintf(format, args...))
}

func PanicId(ctx context.Context, args ...interface{}) { // nolint
	logger.Panic(requestid.String(ctx) + " " + fmt.Sprintln(args...))
}

func PanicfId(ctx context.Context, format string, args ...interface{}) { // nolint
	logger.Panic(requestid.String(ctx) + " " + fmt.Sprintf(format, args...))
}
