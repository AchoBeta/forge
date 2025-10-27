package entity

import (
	"context"
	"forge/pkg/log/zlog"
	"go.uber.org/zap"
)

type User struct {
	UserID   string
	Name     string
	Password string
	Dogs     []*Dog
	// ... ex
}

type Dog struct {
	DogID   string
	DogName string
}

type userCtxKey struct{}

func WithUser(ctx context.Context, user *User) context.Context {
	// 设置用户链路 todo可以考虑在jwt层设置
	ctx = zlog.WithLogKey(ctx, zap.String("user_id", user.UserID))
	ctx = context.WithValue(ctx, userCtxKey{}, *user)
	return ctx
}

func GetUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userCtxKey{}).(*User)
	return user, ok
}
