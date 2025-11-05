package router

import (
	"fmt"
	"forge/biz/types"
	"forge/infra/configs"
	"forge/interface/middleware"
	"forge/pkg/log/zlog"
	"forge/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

var (
	jwtAuthMiddleware gin.HandlerFunc
)

// InitJWTAuth 初始化JWT鉴权中间件
func InitJWTAuth(userService types.IUserService) {
	jwtConfig := configs.Config().GetJWTConfig()

	// 如果secret_key为空，使用默认值（仅开发环境）
	secretKey := jwtConfig.SecretKey
	if secretKey == "" {
		secretKey = "default-secret-key-change-in-production"
		zlog.Warnf("JWT secret_key is empty, using default key. Please set it in config.yaml")
	}

	jwtUtil := util.NewJWTUtil(secretKey, jwtConfig.ExpireHours)
	jwtAuthMiddleware = middleware.JWTAuth(jwtUtil, userService)
}

func RunServer() {
	r := register()
	run(r)
}

func register() (router *gin.Engine) {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.RouterGroup = *r.Group("/api/biz/v1", middleware.AddTracer())

	// 用户服务：不需要JWT的路由（登录、注册、发送验证码、重置密码）
	userGroup := r.Group("user")
	loadUserService(userGroup)

	// 用户服务：需要JWT鉴权的路由（更新头像等）
	userAuthGroup := r.Group("user", jwtAuthMiddleware)
	loadUserAuthService(userAuthGroup)

	// mindmap路由组需要JWT鉴权
	mindMapGroup := r.Group("mindmap", jwtAuthMiddleware)
	loadMindMapService(mindMapGroup)

	// cos路由组需要JWT鉴权
	cosGroup := r.Group("cos", jwtAuthMiddleware)
	loadCOSService(cosGroup)

	return r
}

func run(router *gin.Engine) {
	prot := cast.ToString(configs.Config().GetAppConfig().Port)
	host := configs.Config().GetAppConfig().Host

	zlog.Infof("server run success")
	router.Run(fmt.Sprintf("%s:%s", host, prot))
	zlog.Infof("close run success")
}

const (
	POST   = "POST"
	GET    = "GET"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func loadUserService(r *gin.RouterGroup) {
	r.Handle(POST, "login", Login())

	// 注册接口 user/api/biz/v1/register
	// [POST] /api/biz/v1/user/register
	r.Handle(POST, "register", Register())

	// 发送验证码接口
	// [POST] /api/biz/v1/user/send_code
	r.Handle(POST, "send_code", SendCode())

	// 重置密码接口
	// [POST] /api/biz/v1/user/reset_password
	r.Handle(POST, "reset_password", ResetPassword())
}

// loadUserAuthService 加载需要JWT鉴权的用户服务路由
func loadUserAuthService(r *gin.RouterGroup) {
	// 更新头像接口
	// [PUT] /api/biz/v1/user/avatar
	r.Handle(PUT, "avatar", UpdateAvatar())
}

func loadMindMapService(r *gin.RouterGroup) {
	// 创建思维导图
	// [POST] /api/biz/v1/mindmap
	r.Handle(POST, "", CreateMindMap())

	// 获取思维导图详情
	// [GET] /api/biz/v1/mindmap/:id
	r.Handle(GET, ":id", GetMindMap())

	// 获取思维导图列表
	// [GET] /api/biz/v1/mindmap/list
	r.Handle(GET, "list", ListMindMaps())

	// 更新思维导图
	// [PUT] /api/biz/v1/mindmap/:id
	r.Handle(PUT, ":id", UpdateMindMap())

	// 删除思维导图
	// [DELETE] /api/biz/v1/mindmap/:id
	r.Handle(DELETE, ":id", DeleteMindMap())
}

func loadCOSService(r *gin.RouterGroup) {
	// 获取OSS临时凭证
	// [POST] /api/biz/v1/cos/sts/credentials
	r.Handle(POST, "sts/credentials", GetOSSCredentials())
}
