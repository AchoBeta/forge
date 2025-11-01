package router

import (
	"fmt"
	"forge/infra/configs"
	"forge/interface/middleware"
	"forge/pkg/log/zlog"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func RunServer() {
	r := register()
	run(r)
}

func register() (router *gin.Engine) {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.RouterGroup = *r.Group("/api/biz/v1", middleware.AddTracer())
	loadUserService(r.Group("user"))

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

	// 重置密码接口
	// [POST] /api/biz/v1/user/resetpassword
	r.Handle(POST, "resetpassword", ResetPassword())
}
