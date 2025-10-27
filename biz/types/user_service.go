package types

import (
	"context"
	"forge/biz/entity"
)

type IUserService interface {
	Login(ctx context.Context, username, password string) (*entity.User, error)
}
