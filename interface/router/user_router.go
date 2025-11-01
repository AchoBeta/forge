package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"forge/biz/userservice"
	// "forge/constant"
	"forge/interface/def"
	"forge/interface/handler"
	"forge/pkg/log/zlog"

	// "forge/pkg/loop"
	"forge/pkg/response"
	"forge/util"
)

// 这里就是gin框架的相关接触代码
// 因为解耦的缘故，框架层面的更换不会对内部代码造成任何影响
// router 与hander 应该是一个一对一的关系，有可能会有多对一的关系

// mapServiceErrorToMsgCode 根据应用层返回的错误映射到相应的错误码
func mapServiceErrorToMsgCode(err error) response.MsgCode {
	if err == nil {
		return response.SUCCESS
	}

	// 对应 code_der.go
	// 使用 errors.Is 进行哨兵错误匹配，更加健壮  避免通过字符串匹配来判断
	if errors.Is(err, userservice.ErrUserNotFound) {
		return response.USER_ACCOUNT_NOT_EXIST
	}

	if errors.Is(err, userservice.ErrUserAlreadyExists) {
		return response.USER_ACCOUNT_ALREADY_EXIST
	}

	if errors.Is(err, userservice.ErrInvalidParams) {
		return response.PARAM_NOT_VALID
	}

	if errors.Is(err, userservice.ErrPasswordMismatch) {
		return response.USER_PASSWORD_DIFFERENT
	}

	if errors.Is(err, userservice.ErrCredentialsIncorrect) {
		return response.USER_CREDENTIALS_ERROR
	}

	if errors.Is(err, userservice.ErrUnsupportedAccountType) {
		return response.PARAM_NOT_VALID
	}

	if errors.Is(err, userservice.ErrInternalError) {
		return response.INTERNAL_ERROR
	}

	// 密码强度校验错误
	if errors.Is(err, util.ErrPasswordTooShort) {
		return response.PARAM_NOT_VALID
	}
	if errors.Is(err, util.ErrPasswordTooWeak) {
		return response.PARAM_NOT_VALID
	}
	if errors.Is(err, util.ErrPasswordTooLong) {
		return response.PARAM_NOT_VALID
	}

	// 默认返回通用错误
	return response.COMMON_FAIL
}

// Login
//
//	@Description:[POST] /api/biz/v1/user/login
//	@return gin.HandlerFunc
func Login() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		req := &def.LoginReq{}
		ctx := gCtx.Request.Context()

		// 绑定JSON请求体
		if err := gCtx.ShouldBindJSON(req); err != nil {
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.INVALID_PARAMS.Code,
				Message: response.INVALID_PARAMS.Msg,
				Data:    def.LoginResp{Success: false},
			})
			return
		}

		// TODO: cozeloop配置好后启用
		// ctx, sp := loop.GetNewSpan(ctx, "login", constant.LoopSpanType_Root)
		rsp, err := handler.GetHandler().Login(ctx, req)
		// loop.SetSpanAllInOne(ctx, sp, req, rsp, err)
		zlog.CtxAllInOne(ctx, "login", req, rsp, err)

		// 语法糖包装
		r := response.NewResponse(gCtx)
		if err != nil {
			// 这里是随便写的一个错误码，实际错误码会更加复杂，如何设计更加优雅？
			// 一个handler会返回错误码的
			// r.Error(response.USER_NOT_LOGIN) // todo开放性问题，如何借助 errors.wrap 和errors.Is来更优雅返回msgcode
			// 根据错误类型返回更具体的错误码
			msgCode := mapServiceErrorToMsgCode(err)
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    msgCode.Code,
				Message: msgCode.Msg,
				Data:    def.LoginResp{Success: false},
			})
			return
		} else {
			r.Success(rsp)
		}
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
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.INVALID_PARAMS.Code,
				Message: response.INVALID_PARAMS.Msg,
				Data:    def.RegisterResp{Success: false},
			})
			return
		}

		rsp, err := handler.GetHandler().Register(ctx, req)
		r := response.NewResponse(gCtx)

		if err != nil {
			// 根据 err 的类型返回更具体的错误码
			msgCode := mapServiceErrorToMsgCode(err)
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    msgCode.Code,
				Message: msgCode.Msg,
				Data:    def.RegisterResp{Success: false},
			})
			return
		}
		r.Success(rsp)
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
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.INVALID_PARAMS.Code,
				Message: response.INVALID_PARAMS.Msg,
				Data:    def.ResetPasswordResp{Success: false},
			})
			return
		}

		rsp, err := handler.GetHandler().ResetPassword(ctx, req)
		r := response.NewResponse(gCtx)

		if err != nil {
			// 根据服务层返回的错误类型，返回给客户端更精确的错误信息
			msgCode := mapServiceErrorToMsgCode(err)
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    msgCode.Code,
				Message: msgCode.Msg,
				Data:    def.ResetPasswordResp{Success: false},
			})
			return
		}
		r.Success(rsp)
	}
}
