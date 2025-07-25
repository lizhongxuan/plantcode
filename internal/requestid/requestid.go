package requestid

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid" //nolint
	"time"                   //nolint
)

type contextValue struct {
	ID     string
	Create time.Time
}

type contextKey string

const (
	xRequestID contextKey = "X-Request-ID"
)

// NewContext 使用指定的requestid创建context
func NewContext(ctx context.Context, id string, create time.Time) context.Context {
	val := &contextValue{}
	val.ID = id
	val.Create = create
	return context.WithValue(ctx, xRequestID, val)
}

// WithContext 新建一个context，自动分配requestid
func WithContext(ctx context.Context) context.Context {
	val := &contextValue{}
	id := uuid.New()
	val.ID = hex.EncodeToString(id[:])

	val.Create = time.Now()
	return context.WithValue(ctx, xRequestID, val)
}

func String(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(xRequestID); v != nil {
		val := v.(*contextValue)
		return fmt.Sprintf("[reqid:%s cost=%.3f]", val.ID, time.Since(val.Create).Seconds())
	}
	return ""
}

func GetID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(xRequestID); v != nil {
		val := v.(*contextValue)
		return val.ID
	}
	return ""
}

// Cost 获取从requestid创建，到现在所消耗的时间
func Cost(ctx context.Context) time.Duration {
	if v := ctx.Value(xRequestID); v != nil {
		val := v.(*contextValue)
		return time.Since(val.Create)
	}
	return time.Duration(0)
}

// CreateTime 获取从requestid创建的时间
func CreateTime(ctx context.Context) time.Time {
	if v := ctx.Value(xRequestID); v != nil {
		val := v.(*contextValue)
		return val.Create
	}
	return time.Time{}
}
