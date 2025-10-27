package userservice

import (
	"context"
	"forge/biz/adapter"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/infra/coze"
	"forge/pkg/log/zlog"
)

// 最好的设计方案：
// infra的所有函数都是通过接口来用的

type UserServiceImpl struct {
	userRepo    repo.UserRepo
	cozeService adapter.CozeService
}

func NewUserServiceImpl(
	userRepo repo.UserRepo,
	cozeService adapter.CozeService) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:    userRepo,
		cozeService: cozeService,
	}

}

func (u *UserServiceImpl) Login(ctx context.Context, username, password string) (*entity.User, error) {
	// 这里可以看你自己需要是否加函数级打点
	user := &entity.User{
		Name:     username,
		Password: password,
	}
	err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// 这里如果要调用coze service有两种方法
	// 第一种
	result, err := u.cozeService.RunWorkflow(ctx, &adapter.RunWorkflowReq{})
	if err != nil {
		zlog.CtxErrorf(ctx, "run workflow failed: %w", err)
		return nil, err
	}
	zlog.CtxInfof(ctx, "result:%v", result)
	// 第二种
	result, err = coze.GetCozeService().RunWorkflow(ctx, &adapter.RunWorkflowReq{})
	if err != nil {
		zlog.CtxErrorf(ctx, "run workflow failed: %w", err)
		return nil, err
	}
	zlog.CtxInfof(ctx, "result:%v", result)
	// anything you want

	return user, nil
}
