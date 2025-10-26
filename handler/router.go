package handler

import (
	"forge/infra/configs"
	"github.com/gin-gonic/gin"
)

func LoadRouter(config configs.IConfig) (router *gin.Engine) {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.RouterGroup = *r.Group("/api/biz/v1")
	loadUserService(r.Group("user"))

	return router
}

const (
	POST   = "POST"
	GET    = "GET"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func loadUserService(r *gin.RouterGroup) {
	r.Handle(POST, "login")
}

// 使用gin框架写个简单路由把
