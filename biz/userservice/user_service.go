package userservice

import (
	"context"
	"fmt"
	"forge/biz/adapter"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/biz/types"
	"forge/infra/coze"
	"forge/pkg/log/zlog"
	"forge/util"
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
		UserName: username,
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

// Register 基于手机号/邮箱进行注册
func (u *UserServiceImpl) Register(ctx context.Context, req *types.RegisterParams) (*entity.User, error) {
	// 基本校验
	if req.Account == "" || req.AccountType == "" || req.Password == "" {
		zlog.CtxErrorf(ctx, "invalid params for register")
		return nil, fmt.Errorf("invalid params for register")
	}

	// 检查账号是否已存在
	switch req.AccountType {
	case "phone":
		query := repo.NewUserQueryByPhone(req.Account)
		if exist, _ := u.userRepo.GetUser(ctx, query); exist != nil {
			zlog.CtxErrorf(ctx, "phone already registered")
			return nil, fmt.Errorf("phone already registered")
		}
	case "email":
		query := repo.NewUserQueryByEmail(req.Account)
		if exist, _ := u.userRepo.GetUser(ctx, query); exist != nil {
			zlog.CtxErrorf(ctx, "email already registered")
			return nil, fmt.Errorf("email already registered")
		}
	default:
		zlog.CtxErrorf(ctx, "unsupported accountType: %s", req.AccountType)
		return nil, fmt.Errorf("unsupported accountType: %s", req.AccountType)
	}

	// 校验验证码 code（短信/邮箱） 占位

	//------------------------------------------------

	// 生成用户ID  snowflake雪花id
	userID, err := util.GenerateStringID()
	if err != nil {
		zlog.CtxErrorf(ctx, "generate user id failed: %v", err)
		return nil, fmt.Errorf("generate user id failed: %w", err)
	}
	//

	// 加密密码
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 组装实体 仓储接口写入数据库持久化
	user := &entity.User{
		UserID:   userID,
		UserName: req.UserName,
		Password: hash,
		// 根据 accountType 填写登录方式字段
	}
	if req.AccountType == "phone" {
		user.Phone = req.Account
		user.PhoneVerified = true
	} else if req.AccountType == "email" {
		user.Email = req.Account
		user.EmailVerified = true
	}

	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ResetPassword 重置密码
func (u *UserServiceImpl) ResetPassword(ctx context.Context, req *types.ResetPasswordParams) error {
	// 参数校验
	if req == nil {
		zlog.CtxErrorf(ctx, "reset password request is nil")
		return fmt.Errorf("invalid params for reset password")
	}
	if req.Account == "" || req.AccountType == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		zlog.CtxErrorf(ctx, "invalid params for reset password: missing required fields")
		return fmt.Errorf("invalid params for reset password: missing required fields")
	}

	// 校验两次密码一致性
	if req.NewPassword != req.ConfirmPassword {
		zlog.CtxErrorf(ctx, "password and confirm password do not match")
		return fmt.Errorf("password and confirm password do not match")
	}

	// 根据账号类型查找用户
	var user *entity.User
	var err error
	switch req.AccountType {
	case "phone":
		query := repo.NewUserQueryByPhone(req.Account)
		user, err = u.userRepo.GetUser(ctx, query)
		if err != nil {
			zlog.CtxErrorf(ctx, "get user by phone failed: %v", err)
			return fmt.Errorf("user not found")
		}
		if user == nil {
			zlog.CtxErrorf(ctx, "user not found by phone: %s", req.Account)
			return fmt.Errorf("user not found")
		}
	case "email":
		query := repo.NewUserQueryByEmail(req.Account)
		user, err = u.userRepo.GetUser(ctx, query)
		if err != nil {
			zlog.CtxErrorf(ctx, "get user by email failed: %v", err)
			return fmt.Errorf("user not found")
		}
		if user == nil {
			zlog.CtxErrorf(ctx, "user not found by email: %s", req.Account)
			return fmt.Errorf("user not found")
		}
	default:
		zlog.CtxErrorf(ctx, "unsupported accountType: %s", req.AccountType)
		return fmt.Errorf("unsupported accountType: %s", req.AccountType)
	}

	// 4. 校验验证码 code（短信/邮箱），此处预留

	// 验证码校验逻辑应该在这里实现

	// 加密新密码
	hash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		zlog.CtxErrorf(ctx, "hash password failed: %v", err)
		return fmt.Errorf("failed to hash password")
	}

	// 更新用户密码
	password := hash
	updateInfo := &repo.UserUpdateInfo{
		UserID:   user.UserID,
		Password: &password,
	}
	if err := u.userRepo.UpdateUser(ctx, updateInfo); err != nil {
		zlog.CtxErrorf(ctx, "update password failed: %v", err)
		return fmt.Errorf("failed to update password")
	}

	zlog.CtxInfof(ctx, "reset password successfully for user: %s", user.UserID)
	return nil
}
