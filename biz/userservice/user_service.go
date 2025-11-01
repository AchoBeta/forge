package userservice

import (
	"context"
	"errors"
	"fmt"
	"forge/biz/adapter"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/biz/types"
	"forge/pkg/log/zlog"
	"forge/util"
)

var (
	// ErrUserNotFound 表示用户不存在
	ErrUserNotFound = errors.New("user not found")
)

// 最好的设计方案：
// infra的所有函数都是通过接口来用的

type UserServiceImpl struct {
	userRepo    repo.UserRepo
	cozeService adapter.CozeService
	jwtUtil     *util.JWTUtil
}

// 默认JWT配置（可在配置文件中配置）
const (
	DefaultJWTSecretKey   = "forge-secret-key-change-in-production"
	DefaultJWTExpireHours = 24 * 7 // 过期时间7天
)

func NewUserServiceImpl(
	userRepo repo.UserRepo,
	cozeService adapter.CozeService) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:    userRepo,
		cozeService: cozeService,
		jwtUtil:     util.NewJWTUtil(DefaultJWTSecretKey, DefaultJWTExpireHours),
	}
}

// Login 登录：根据账号和密码进行登录
func (u *UserServiceImpl) Login(ctx context.Context, account, accountType, password string) (*entity.User, string, error) {
	// 参数校验
	if account == "" || accountType == "" || password == "" {
		zlog.CtxErrorf(ctx, "invalid params for login: account, accountType or password is empty")
		return nil, "", fmt.Errorf("invalid params for login")
	}

	// 根据账号类型查找用户
	user, err := u.findUserByAccount(ctx, account, accountType)
	if err != nil {
		// 如果用户不存在，返回错误
		if errors.Is(err, ErrUserNotFound) {
			zlog.CtxErrorf(ctx, "user not found: %s", account)
			return nil, "", fmt.Errorf("account or password incorrect")
		}
		// 其他错误（数据库错误等）
		return nil, "", err
	}

	// 验证密码
	match, err := util.ComparePassword(user.Password, password)
	if err != nil {
		zlog.CtxErrorf(ctx, "compare password failed: %v", err)
		return nil, "", fmt.Errorf("internal error: failed to verify password")
	}
	if !match {
		zlog.CtxErrorf(ctx, "password incorrect for user: %s", user.UserID)
		return nil, "", fmt.Errorf("account or password incorrect")
	}

	// 生成JWT token
	token, err := u.jwtUtil.GenerateToken(user.UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "generate token failed: %v", err)
		return nil, "", fmt.Errorf("failed to generate token")
	}

	// 方法一  通过注入的 cozeService 接口调用 哇哦
	result, err := u.cozeService.RunWorkflow(ctx, &adapter.RunWorkflowReq{})
	if err != nil {
		zlog.CtxErrorf(ctx, "run workflow failed: %v", err)
	} else {
		zlog.CtxInfof(ctx, "result:%v", result)
	}

	// 方法二
	// result, err = coze.GetCozeService().RunWorkflow(ctx, &adapter.RunWorkflowReq{})
	// if err != nil {
	// 	zlog.CtxErrorf(ctx, "run workflow failed: %v", err)
	// 	return nil, "", err
	// }
	// zlog.CtxInfof(ctx, "result:%v", result)
	// ============================================================

	// 更新最后登录时间（可选）
	// lastLoginAt := time.Now()
	// updateInfo := &repo.UserUpdateInfo{
	// 	UserID:     user.UserID,
	// 	LastLoginAt: &lastLoginAt,
	// }
	// _ = u.userRepo.UpdateUser(ctx, updateInfo)

	zlog.CtxInfof(ctx, "login success for user: %s", user.UserID)
	return user, token, nil
}

// Register 基于手机号/邮箱进行注册
func (u *UserServiceImpl) Register(ctx context.Context, req *types.RegisterParams) (*entity.User, error) {
	// 基本校验
	if req.Account == "" || req.AccountType == "" || req.Password == "" {
		zlog.CtxErrorf(ctx, "invalid params for register")
		return nil, fmt.Errorf("invalid params for register")
	}

	// 检查账号是否已存在
	existUser, err := u.findUserByAccount(ctx, req.Account, req.AccountType)
	if err != nil {
		// 账号不存在，可以继续注册
		if errors.Is(err, ErrUserNotFound) {
			// 用户不存在，继续注册流程
		} else {
			// 其他错误，直接返回
			return nil, err
		}
	} else if existUser != nil {
		// 用户已存在，返回错误
		var accountField string
		if req.AccountType == types.AccountTypePhone {
			accountField = "phone"
		} else {
			accountField = "email"
		}
		zlog.CtxErrorf(ctx, "%s already registered: %s", accountField, req.Account)
		return nil, fmt.Errorf("%s already registered", accountField)
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
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 组装实体 仓储接口写入数据库持久化
	user := &entity.User{
		UserID:   userID,
		UserName: req.UserName,
		Password: hash,
		// 根据 accountType 填写登录方式字段
	}
	if req.AccountType == types.AccountTypePhone {
		user.Phone = req.Account
		user.PhoneVerified = true
	} else if req.AccountType == types.AccountTypeEmail {
		user.Email = req.Account
		user.EmailVerified = true
	}

	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// findUserByAccount 根据账号类型查找用户 抽离重复判断逻辑
// 返回值说明：
//   - 如果返回错误不为nil，表示数据库查询出错（内部错误）或账号类型不支持
//   - 如果用户为nil且错误为nil，表示用户不存在，返回"user not found"错误
//   - 如果用户不为nil，表示找到用户，正常返回
func (u *UserServiceImpl) findUserByAccount(ctx context.Context, account, accountType string) (*entity.User, error) {
	var query repo.UserQuery
	var accountField string

	switch accountType {
	case types.AccountTypePhone:
		query = repo.NewUserQueryByPhone(account)
		accountField = "phone"
	case types.AccountTypeEmail:
		query = repo.NewUserQueryByEmail(account)
		accountField = "email"
	default:
		zlog.CtxErrorf(ctx, "unsupported accountType: %s", accountType)
		return nil, fmt.Errorf("unsupported accountType: %s", accountType)
	}

	user, err := u.userRepo.GetUser(ctx, query)
	if err != nil {
		// 数据库查询错误，返回内部错误
		zlog.CtxErrorf(ctx, "failed to get user by %s: %v", accountField, err)
		return nil, fmt.Errorf("internal error: failed to query user")
	}

	if user == nil {
		// 用户不存在
		zlog.CtxErrorf(ctx, "user not found by %s: %s", accountField, account)
		return nil, ErrUserNotFound
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
	user, err := u.findUserByAccount(ctx, req.Account, req.AccountType)
	if err != nil {
		return err
	}

	// 4. 校验验证码 code（短信/邮箱），此处预留

	// 验证码校验逻辑应该在这里实现

	// 加密新密码
	hash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		zlog.CtxErrorf(ctx, "hash password failed: %v", err)
		return fmt.Errorf("failed to hash password: %w", err)
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
