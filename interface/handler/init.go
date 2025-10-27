package handler

import (
	"context"
	"forge/biz/types"
	"forge/interface/def"
)

type IHandler interface {
	Login(ctx context.Context, req *def.LoginReq) (rsp *def.LoginResp, err error)
}

var handler IHandler

type Handler struct {
	UserService types.IUserService
}

func GetHandler() IHandler {
	return handler
}
func MustInitHandler(userService types.IUserService) {
	err := InitHandler(userService)
	if err != nil {
		panic(err)
	}
}

func InitHandler(userService types.IUserService) error {
	handler = &Handler{
		UserService: userService,
	}
	return nil
}
