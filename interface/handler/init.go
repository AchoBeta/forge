package handler

import (
	"context"
	"forge/biz/types"
	"forge/interface/def"
)

type IHandler interface {
	Login(ctx context.Context, req *def.LoginReq) (rsp *def.LoginResp, err error)
	// Register: 注册 暂无第三方
	Register(ctx context.Context, req *def.RegisterReq) (rsp *def.RegisterResp, err error)
	// ResetPassword: 重置密码
	ResetPassword(ctx context.Context, req *def.ResetPasswordReq) (rsp *def.ResetPasswordResp, err error)
	// SendCode: 发送验证码  ！邮件！
	SendCode(ctx context.Context, req *def.SendVerificationCodeReq) (rsp *def.SendVerificationCodeResp, err error)

	// MindMap: 思维导图相关接口
	CreateMindMap(ctx context.Context, req *def.CreateMindMapReq) (rsp *def.CreateMindMapResp, err error)
	GetMindMap(ctx context.Context, mapID string) (rsp *def.GetMindMapResp, err error)
	ListMindMaps(ctx context.Context, req *def.ListMindMapsReq) (rsp *def.ListMindMapsResp, err error)
	UpdateMindMap(ctx context.Context, mapID string, req *def.UpdateMindMapReq) (rsp *def.UpdateMindMapResp, err error)
	DeleteMindMap(ctx context.Context, mapID string) (rsp *def.DeleteMindMapResp, err error)
}

var handler IHandler

type Handler struct {
	UserService    types.IUserService
	MindMapService types.IMindMapService
}

func GetHandler() IHandler {
	return handler
}
func MustInitHandler(userService types.IUserService, mindMapService types.IMindMapService) {
	err := InitHandler(userService, mindMapService)
	if err != nil {
		panic(err)
	}
}

func InitHandler(userService types.IUserService, mindMapService types.IMindMapService) error {
	handler = &Handler{
		UserService:    userService,
		MindMapService: mindMapService,
	}
	return nil
}
