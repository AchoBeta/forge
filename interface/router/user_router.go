package router

import (
	"forge/constant"
	"forge/interface/def"
	"forge/interface/handler"
	"forge/pkg/log/zlog"
	"forge/pkg/loop"
	"forge/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 这里就是gin框架的相关接触代码
// 因为解耦的缘故，框架层面的更换不会对内部代码造成任何影响
// router 与hander 应该是一个一对一的关系，有可能会有多对一的关系

// Login
//
//	@Description:[POST] /api/biz/v1/user/login
//	@return gin.HandlerFunc
func Login() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		// 这里以某种方式获组装到了请求体
		req := &def.LoginReq{}
		ctx := gCtx.Request.Context()
		ctx, sp := loop.GetNewSpan(ctx, "login", constant.LoopSpanType_Root)
		rsp, err := handler.GetHandler().Login(ctx, req)
		loop.SetSpanAllInOne(ctx, sp, req, rsp, err)
		zlog.CtxAllInOne(ctx, "login", req, rsp, err)

		// 语法糖包装
		r := response.NewResponse(gCtx)
		if err != nil {
			// 这里是随便写的一个错误码，实际错误码会更加复杂，如何设计更加优雅？
			// 一个handler会返回错误码的
			r.Error(response.USER_NOT_LOGIN) // todo开放性问题，如何借助 errors.wrap 和errors.Is来更优雅返回msgcode
		} else {
			r.Success(rsp)
		}
		return
	}
}

// Register
//
//	@Description:[POST] /api/biz/v1/user/register
//	@return gin.HandlerFunc
func Register() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		req := &def.RegisterReq{}
		// 统一从 gin 上下文取出 request 的 context，供后续业务调用使用
		ctx := gCtx.Request.Context()
		if err := gCtx.ShouldBindJSON(req); err != nil {
			r := response.NewResponse(gCtx)
			r.Error(response.USER_NOT_LOGIN) //USER_NOT_LOGIN  = MsgCode{Code: 2001, Msg: "用户未登录"}
			return
		}

		rsp, err := handler.GetHandler().Register(ctx, req)

		if err != nil {
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{Code: response.USER_NOT_LOGIN.Code, Message: response.USER_NOT_LOGIN.Msg, Data: struct{}{}})
			return
		}
		gCtx.JSON(http.StatusOK, response.JsonMsgResult{Code: response.SUCCESS_CODE, Message: response.SUCCESS_MSG, Data: rsp})
	}
}

// ResetPassword
//
//	@Description:[POST] /api/biz/v1/user/reset_password
//	@return gin.HandlerFunc
func ResetPassword() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		req := &def.ResetPasswordReq{}
		// 统一从 gin 上下文取出 request 的 context
		ctx := gCtx.Request.Context()
		if err := gCtx.ShouldBindJSON(req); err != nil {
			r := response.NewResponse(gCtx)
			r.Error(response.USER_NOT_LOGIN)
			return
		}

		rsp, err := handler.GetHandler().ResetPassword(ctx, req)

		if err != nil {
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{Code: response.USER_NOT_LOGIN.Code, Message: response.USER_NOT_LOGIN.Msg, Data: struct{}{}})
			return
		}
		gCtx.JSON(http.StatusOK, response.JsonMsgResult{Code: response.SUCCESS_CODE, Message: response.SUCCESS_MSG, Data: rsp})
	}
}
