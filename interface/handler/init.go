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

	// COS: OSS凭证相关接口
	GetOSSCredentials(ctx context.Context, req *def.GetOSSCredentialsReq) (rsp *def.GetOSSCredentialsResp, err error)

	//AiChat: ai对话相关
	SendMessage(ctx context.Context, req *def.ProcessUserMessageRequest) (*def.ProcessUserMessageResponse, error)
	SaveNewConversation(ctx context.Context, req *def.SaveNewConversationRequest) (*def.SaveNewConversationResponse, error)
	GetConversationList(ctx context.Context, req *def.GetConversationListRequest) (*def.GetConversationListResponse, error)
	DelConversation(ctx context.Context, req *def.DelConversationRequest) (*def.DelConversationResponse, error)
	GetConversation(ctx context.Context, req *def.GetConversationRequest) (*def.GetConversationResponse, error)
}

var handler IHandler

type Handler struct {
	UserService    types.IUserService
	MindMapService types.IMindMapService
	COSService     types.ICOSService
	AiChatService  types.IAiChatService
}

func GetHandler() IHandler {
	return handler
}
func MustInitHandler(userService types.IUserService, mindMapService types.IMindMapService, cosService types.ICOSService, aiChatService types.IAiChatService) {
	err := InitHandler(userService, mindMapService, cosService, aiChatService)
	if err != nil {
		panic(err)
	}
}

func InitHandler(userService types.IUserService, mindMapService types.IMindMapService, cosService types.ICOSService, aiChatService types.IAiChatService) error {
	handler = &Handler{
		UserService:    userService,
		MindMapService: mindMapService,
		COSService:     cosService,
		AiChatService:  aiChatService,
	}
	return nil
}
