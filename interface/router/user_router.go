package router

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"forge/biz/cosservice"
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

	// 验证码错误
	if errors.Is(err, userservice.ErrVerificationCodeIncorrect) {
		return response.CAPTCHA_ERROR
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

	// COS相关错误
	if errors.Is(err, cosservice.ErrInvalidParams) {
		return response.PARAM_NOT_VALID
	}
	if errors.Is(err, cosservice.ErrPermissionDenied) {
		return response.INSUFFICENT_PERMISSIONS
	}
	if errors.Is(err, cosservice.ErrInternalError) {
		return response.INTERNAL_FILE_UPLOAD_ERROR
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

// SendCode
//
//	@Description:[POST] /api/biz/v1/user/send_code
//	@return gin.HandlerFunc
func SendCode() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		req := &def.SendVerificationCodeReq{}
		// 统一从 gin 上下文取出 request 的 context
		ctx := gCtx.Request.Context()
		if err := gCtx.ShouldBindJSON(req); err != nil {
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.INVALID_PARAMS.Code,
				Message: response.INVALID_PARAMS.Msg,
				Data:    def.SendVerificationCodeResp{Success: false},
			})
			return
		}

		rsp, err := handler.GetHandler().SendCode(ctx, req)
		r := response.NewResponse(gCtx)

		if err != nil {
			// 根据服务层返回的错误类型，返回给客户端更精确的错误信息
			msgCode := mapServiceErrorToMsgCode(err)
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    msgCode.Code,
				Message: msgCode.Msg,
				Data:    def.SendVerificationCodeResp{Success: false},
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

// UpdateAvatar
//
//	@Description:[POST] /api/biz/v1/user/avatar
//	@return gin.HandlerFunc
func UpdateAvatar() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		ctx := gCtx.Request.Context()

		// 设置文件大小限制（8MB）
		gCtx.Request.ParseMultipartForm(8 << 20)

		// 接收文件
		file, err := gCtx.FormFile("avatar") // "avatar" 是前端表单字段名
		if err != nil {
			// 检查是否是文件大小错误
			if strings.Contains(err.Error(), "too large") || strings.Contains(err.Error(), "request body too large") {
				zlog.CtxErrorf(ctx, "file too large: %v", err)
				gCtx.JSON(http.StatusOK, response.JsonMsgResult{
					Code:    response.PARAM_FILE_SIZE_TOO_BIG.Code,
					Message: response.PARAM_FILE_SIZE_TOO_BIG.Msg,
					Data:    def.UpdateAvatarResp{Success: false},
				})
				return
			}
			zlog.CtxErrorf(ctx, "failed to get file from form: %v", err)
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.PARAM_NOT_VALID.Code,
				Message: response.PARAM_NOT_VALID.Msg,
				Data:    def.UpdateAvatarResp{Success: false},
			})
			return
		}

		// 检查文件大小
		if file.Size > 5*1024*1024 { // 5MB
			zlog.CtxErrorf(ctx, "file size too large: %d bytes", file.Size)
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.PARAM_FILE_SIZE_TOO_BIG.Code,
				Message: response.PARAM_FILE_SIZE_TOO_BIG.Msg,
				Data:    def.UpdateAvatarResp{Success: false},
			})
			return
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			zlog.CtxErrorf(ctx, "failed to open file: %v", err)
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.INTERNAL_FILE_UPLOAD_ERROR.Code,
				Message: response.INTERNAL_FILE_UPLOAD_ERROR.Msg,
				Data:    def.UpdateAvatarResp{Success: false},
			})
			return
		}
		defer src.Close() // 确保关闭

		// 读取文件内容
		fileData, err := io.ReadAll(src)
		if err != nil {
			zlog.CtxErrorf(ctx, "failed to read file: %v", err)
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    response.INTERNAL_FILE_UPLOAD_ERROR.Code,
				Message: response.INTERNAL_FILE_UPLOAD_ERROR.Msg,
				Data:    def.UpdateAvatarResp{Success: false},
			})
			return
		}

		// 构建请求对象
		req := &def.UpdateAvatarReq{
			FileData: fileData,
			Filename: file.Filename,
		}

		// 调用handler
		rsp, err := handler.GetHandler().UpdateAvatar(ctx, req)
		r := response.NewResponse(gCtx)

		if err != nil {
			msgCode := mapServiceErrorToMsgCode(err)
			gCtx.JSON(http.StatusOK, response.JsonMsgResult{
				Code:    msgCode.Code,
				Message: msgCode.Msg,
				Data:    def.UpdateAvatarResp{Success: false},
			})
			return
		}
		r.Success(rsp)
	}
}
