package handler

import (
	"context"
	"forge/constant"
	"forge/interface/caster"
	"forge/interface/def"
	"forge/pkg/log/zlog"
	"forge/pkg/loop"
)

func (h *Handler) Login(ctx context.Context, req *def.LoginReq) (rsp *def.LoginResp, err error) {

	// 这里用作handler级别的链路追踪
	ctx, sp := loop.GetNewSpan(ctx, "handler.login", constant.LoopSpanType_Handle)
	defer func() {
		zlog.CtxAllInOne(ctx, "handler.login", req, rsp, err)
		loop.SetSpanAllInOne(ctx, sp, req, rsp, err)
	}()

	// 这里可能会做更复杂的service编排
	// 为什么我们会有service和handler的区分？
	// 我的理解是，service我们更倾向于做一个原子能力，比如某个动作
	// 但实际业务可能需要一次接口请求先做a再做b再做c，再返回结果
	// 所以这里这么做区分
	// 同时，发布事件应该也在handler层做，service层做就会腐化（引入与你无关的代码）
	user, err := h.UserService.Login(ctx, req.UserName, req.Password)

	// 并且如果两个service有交集：一个service要调用另外一个service的代码
	// 这时候你就要考虑你的两个service是否要合并了
	// 一个好的service他天然就与其他service独立的。
	// 如果确实有这个场景，也可以引入防腐层

	// 注册完成后可能还会有其他步骤，比如送积分啊，之类的。
	// 如何鉴别是在handler里面调用service还是说注册订阅事件消费呢？
	// 最简单的做法是看你的这个操作是强依赖还是弱依赖的
	// 强依赖代表如果接下来这个操作失败了，你的整个接口返回会因此有某些改变
	// 弱依赖是指如果接下来这个操作失败了，不应该影响你整个接口的完成情况，
	// 弱依赖最优雅是 1. 本地数据总线 2. 消息队列
	userDto := caster.CastUserDO2DTO(user)
	rsp.User = userDto
	return
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

	rsp = &def.RegisterResp{}
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
