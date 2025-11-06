package handler

import (
	"context"

	// "forge/constant"
	"forge/biz/entity"
	"forge/biz/userservice"
	"forge/interface/caster"
	"forge/interface/def"
	"forge/pkg/log/zlog"
	// "forge/pkg/loop"
)

func (h *Handler) Login(ctx context.Context, req *def.LoginReq) (rsp *def.LoginResp, err error) {

	// 这里用作handler级别的链路追踪 - TODO: cozeloop配置好后启用
	// ctx, sp := loop.GetNewSpan(ctx, "handler.login", constant.LoopSpanType_Handle)
	defer func() {
		zlog.CtxAllInOne(ctx, "handler.login", req, rsp, err)
		// loop.SetSpanAllInOne(ctx, sp, req, rsp, err)
	}()

	// 这里可能会做更复杂的service编排
	// 为什么我们会有service和handler的区分？
	// 我的理解是，service我们更倾向于做一个原子能力，比如某个动作
	// 但实际业务可能需要一次接口请求先做a再做b再做c，再返回结果
	// 所以这里这么做区分
	// 同时，发布事件应该也在handler层做，service层做就会腐化（引入与你无关的代码）
	// 调用服务层登录
	user, token, err := h.UserService.Login(ctx, req.Account, req.AccountType, req.Password)
	if err != nil {
		return nil, err
	}

	// 组装响应
	rsp = &def.LoginResp{
		Token:    token,
		UserID:   user.UserID,
		UserName: user.UserName,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		Email:    user.Email,
		Success:  true, // 登录成功
	}
	return rsp, nil
}

func (h *Handler) Register(ctx context.Context, req *def.RegisterReq) (rsp *def.RegisterResp, err error) {
	//

	// DTO -> Service 层表单
	params := caster.CastRegisterReq2Params(req)

	// 请求限流、验证验证码等 占位

	// -------------------------

	// 向下调用服务层
	_, err = h.UserService.Register(ctx, params)
	if err != nil {
		return nil, err
	}

	rsp = &def.RegisterResp{
		Success: true,
	}
	return rsp, nil
}

func (h *Handler) ResetPassword(ctx context.Context, req *def.ResetPasswordReq) (rsp *def.ResetPasswordResp, err error) {
	// DTO -> Service 层表单
	params := caster.CastResetPasswordReq2Params(req)

	// 向下调用服务层 重置密码函数
	err = h.UserService.ResetPassword(ctx, params)
	if err != nil {
		return nil, err
	}

	rsp = &def.ResetPasswordResp{
		Success: true,
	}
	return rsp, nil
}

func (h *Handler) SendCode(ctx context.Context, req *def.SendVerificationCodeReq) (rsp *def.SendVerificationCodeResp, err error) {
	// 调用服务层发送验证码
	err = h.UserService.SendVerificationCode(ctx, req.Account, req.AccountType, req.Purpose)
	if err != nil {
		return nil, err
	}

	rsp = &def.SendVerificationCodeResp{
		Success: true,
	}
	return rsp, nil
}

func (h *Handler) GetHome(ctx context.Context) (rsp *def.GetHomeResp, err error) {
	defer func() {
		zlog.CtxAllInOne(ctx, "handler.get_home", nil, rsp, err)
	}()

	// 从context获取当前用户（JWT中间件已注入）
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "user not found in context, this should not happen if JWT middleware works correctly")
		return nil, userservice.ErrPermissionDenied
	}

	// 判断是否有密码
	hasPassword := user.Password != ""

	// 组装响应
	rsp = &def.GetHomeResp{
		UserName:    user.UserName,
		Avatar:      user.Avatar,
		Phone:       user.Phone,
		Email:       user.Email,
		HasPassword: hasPassword,
	}
	return rsp, nil
}

func (h *Handler) UpdateAccount(ctx context.Context, req *def.UpdateAccountReq) (rsp *def.UpdateAccountResp, err error) {
	defer func() {
		zlog.CtxAllInOne(ctx, "handler.update_account", req, rsp, err)
	}()

	// DTO -> Service 层参数转换
	params := caster.CastUpdateAccountReq2Params(req)

	// 调用服务层更新联系方式
	account, err := h.UserService.UpdateAccount(ctx, params)
	if err != nil {
		return nil, err
	}

	rsp = &def.UpdateAccountResp{
		Success: true,
		Account: account,
	}
	return rsp, nil
}
