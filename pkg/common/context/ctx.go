package context

import (
	"context"
	"orderin-server/pkg/common/constant"
)

var mapper = []string{constant.RequestId, constant.OpUserID}

func SetOpUserID(ctx context.Context, opUserID string) context.Context {
	return context.WithValue(ctx, constant.OpUserID, opUserID)
}

func SetRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, constant.RequestId, requestId)
}

func GetOpUserID(ctx context.Context) int64 {
	if ctx.Value(constant.OpUserID) != nil {
		s, ok := ctx.Value(constant.OpUserID).(int64)
		if ok {
			return s
		}
	}
	return 0
}

func GetSuperAdmin(ctx context.Context) bool {
	if ctx.Value(constant.SuperAdmin) != nil {
		s, ok := ctx.Value(constant.SuperAdmin).(bool)
		if ok {
			return s
		}
	}
	return false
}

func GetRequestId(ctx context.Context) string {
	if ctx.Value(constant.RequestId) != nil {
		s, ok := ctx.Value(constant.RequestId).(string)
		if ok {
			return s
		}
	}
	return ""
}

func GetRemoteAddr(ctx context.Context) string {
	if ctx.Value(constant.RemoteAddr) != "" {
		s, ok := ctx.Value(constant.RemoteAddr).(string)
		if ok {
			return s
		}
	}
	return ""
}
