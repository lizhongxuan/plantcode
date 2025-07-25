package log

import (
	"ai-dev-platform/internal/requestid"
	"context"
	"testing"
	"time"
)

func TestRequestIDLog(t *testing.T) {

	ctx := requestid.NewContext(context.Background(), "123", time.Now())
	DebugId(ctx, "debug")
	DebugfId(ctx, "debug %s", "requestid")
	InfoId(ctx, "info")
	InfofId(ctx, "info %s", "requestid")
	WarnId(ctx, "warn")
	WarnfId(ctx, "warn %s", "requestid")
	ErrorId(ctx, "error")
	ErrorfId(ctx, "error %s", "requestid")
}
