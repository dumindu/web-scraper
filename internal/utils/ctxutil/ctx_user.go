package ctxutil

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID    *uuid.UUID
	Email string
}

func SetUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, keyUser, user)
}

func UserFromCtx(ctx context.Context) User {
	user, _ := ctx.Value(keyUser).(User)

	return user
}
