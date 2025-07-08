package ctxutil

import "context"

const (
	keyRequestID key = "requestID"
	keyUser      key = "user"
)

type key string

func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(keyRequestID).(string)

	return requestID
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, keyRequestID, requestID)
}
