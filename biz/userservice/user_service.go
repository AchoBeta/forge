package userservice

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"forge/biz/adapter"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/biz/types"
	"forge/infra/cache"
	"forge/pkg/log/zlog"
	"forge/util"
)

var (
	// ErrUserNotFound 表示用户不存在
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists 表示账号已存在
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidParams 表示参数无效
	ErrInvalidParams = errors.New("invalid params")
	// ErrPasswordMismatch 表示密码不一致
	ErrPasswordMismatch = errors.New("password mismatch")
	// ErrCredentialsIncorrect 表示账号或密码错误
	ErrCredentialsIncorrect = errors.New("credentials incorrect")
	// ErrUnsupportedAccountType 表示不支持的账号类型
	ErrUnsupportedAccountType = errors.New("unsupported account type")
	// ErrInternalError 表示内部错误
	ErrInternalError = errors.New("internal error")
	// ErrPermissionDenied 表示权限被拒绝
	ErrPermissionDenied = errors.New("permission denied")
	// ErrVerificationCodeIncorrect 表示验证码错误
	ErrVerificationCodeIncorrect = errors.New("verification code incorrect")
)

// 最好的设计方案：
// infra的所有函数都是通过接口来用的

type UserServiceImpl struct {
	userRepo     repo.UserRepo
	cozeService  adapter.CozeService
	jwtUtil      *util.JWTUtil
	emailService adapter.EmailService
}

func NewUserServiceImpl(
	userRepo repo.UserRepo,
	cozeService adapter.CozeService,
	jwtUtil *util.JWTUtil,
	emailService adapter.EmailService) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:     userRepo,
		cozeService:  cozeService,
		jwtUtil:      jwtUtil,
		emailService: emailService,
	}
}

// Login 登录：根据账号和密码进行登录
func (u *UserServiceImpl) Login(ctx context.Context, account, accountType, password string) (*entity.User, string, error) {
	// 参数校验
	if account == "" || accountType == "" || password == "" {
		zlog.CtxErrorf(ctx, "invalid params for login: account, accountType or password is empty")
		return nil, "", ErrInvalidParams
	}

	// 根据账号类型查找用户
	user, err := u.findUserByAccount(ctx, account, accountType)
	if err != nil {
		// 如果用户不存在，返回错误
		if errors.Is(err, ErrUserNotFound) {
			zlog.CtxErrorf(ctx, "user not found: %s", account)
			return nil, "", ErrCredentialsIncorrect
		}
		// 其他错误（数据库错误等）
		return nil, "", err
	}

	// 验证密码
	match, err := util.ComparePassword(user.Password, password)
	if err != nil {
		zlog.CtxErrorf(ctx, "compare password failed: %v", err)
		return nil, "", ErrInternalError
	}
	if !match {
		zlog.CtxErrorf(ctx, "password incorrect for user: %s", user.UserID)
		return nil, "", ErrCredentialsIncorrect
	}

	// 生成JWT token
	token, err := u.jwtUtil.GenerateToken(user.UserID)
	if err != nil {
		zlog.CtxErrorf(ctx, "generate token failed: %v", err)
		return nil, "", ErrInternalError
	}

	// 方法一  通过注入的 cozeService 接口调用
	//result, err := u.cozeService.RunWorkflow(ctx, &adapter.RunWorkflowReq{})
	//if err != nil {
	//	zlog.CtxErrorf(ctx, "run workflow failed: %v", err)
	//} else {
	//	zlog.CtxInfof(ctx, "result:%v", result)
	//}

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
		return nil, ErrInvalidParams
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
		return nil, ErrUserAlreadyExists
	}

	// 校验验证码 code（短信/邮箱）
	if err := u.verifyCode(ctx, req.Account, req.AccountType, req.Code, types.PurposeRegister); err != nil {
		return nil, err
	}

	//------------------------------------------------

	// 验证密码强度  按照常规要求设置
	if err := util.ValidatePasswordStrength(req.Password); err != nil {
		zlog.CtxErrorf(ctx, "password strength validation failed: %v", err)
		return nil, err
	}

	// 生成用户ID  snowflake雪花id
	userID, err := util.GenerateStringID()
	if err != nil {
		zlog.CtxErrorf(ctx, "generate user id failed: %v", err)
		return nil, ErrInternalError
	}
	//

	// 加密密码
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, ErrInternalError
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
		return nil, ErrUnsupportedAccountType
	}

	user, err := u.userRepo.GetUser(ctx, query)
	if err != nil {
		// 数据库查询错误，返回内部错误
		zlog.CtxErrorf(ctx, "failed to get user by %s: %v", accountField, err)
		return nil, ErrInternalError
	}

	if user == nil {
		// 用户不存在
		return nil, ErrUserNotFound
	}

	return user, nil
}

// ResetPassword 重置密码
func (u *UserServiceImpl) ResetPassword(ctx context.Context, req *types.ResetPasswordParams) error {
	// 参数校验
	if req == nil {
		zlog.CtxErrorf(ctx, "reset password request is nil")
		return ErrInvalidParams
	}
	if req.Account == "" || req.AccountType == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		zlog.CtxErrorf(ctx, "invalid params for reset password: missing required fields")
		return ErrInvalidParams
	}

	// 校验两次密码一致性
	if req.NewPassword != req.ConfirmPassword {
		zlog.CtxErrorf(ctx, "password and confirm password do not match")
		return ErrPasswordMismatch
	}

	// 根据账号类型查找用户
	user, err := u.findUserByAccount(ctx, req.Account, req.AccountType)
	if err != nil {
		return err
	}

	// 4. 校验验证码 code（短信/邮箱）
	if err := u.verifyCode(ctx, req.Account, req.AccountType, req.Code, types.PurposeResetPassword); err != nil {
		return err
	}

	// 验证新密码强度
	if err := util.ValidatePasswordStrength(req.NewPassword); err != nil {
		zlog.CtxErrorf(ctx, "password strength validation failed: %v", err)
		return err
	}

	// 加密新密码
	hash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		zlog.CtxErrorf(ctx, "hash password failed: %v", err)
		return ErrInternalError
	}

	// 更新用户密码
	password := hash
	updateInfo := &repo.UserUpdateInfo{
		UserID:   user.UserID,
		Password: &password,
	}
	if err := u.userRepo.UpdateUser(ctx, updateInfo); err != nil {
		zlog.CtxErrorf(ctx, "update password failed: %v", err)
		return ErrInternalError
	}

	zlog.CtxInfof(ctx, "reset password successfully for user: %s", user.UserID)
	return nil
}

// GetUserByID 根据用户ID获取用户信息（用于JWT鉴权等场景）
func (u *UserServiceImpl) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	// 参数校验
	if userID == "" {
		zlog.CtxErrorf(ctx, "userID is required")
		return nil, ErrInvalidParams
	}

	// 通过repo查询用户
	query := repo.NewUserQueryByID(userID)
	user, err := u.userRepo.GetUser(ctx, query)
	if err != nil {
		zlog.CtxErrorf(ctx, "failed to get user by ID: %v", err)
		return nil, ErrInternalError
	}

	if user == nil {
		zlog.CtxWarnf(ctx, "user not found: %s", userID)
		return nil, ErrUserNotFound
	}

	// 检查用户状态（业务逻辑应该在service层）
	if user.Status != entity.UserStatusActive {
		zlog.CtxWarnf(ctx, "user is disabled: %s", userID)
		return nil, ErrPermissionDenied
	}

	return user, nil
}

// SendVerificationCode 发送验证码
func (u *UserServiceImpl) SendVerificationCode(ctx context.Context, account, accountType, purpose string) error {
	// 参数校验
	if account == "" || accountType == "" || purpose == "" {
		zlog.CtxErrorf(ctx, "invalid params for send verification code")
		return ErrInvalidParams
	}

	// 目前只支持邮箱验证码
	if accountType != types.AccountTypeEmail {
		zlog.CtxErrorf(ctx, "unsupported account type for verification: %s", accountType)
		return ErrUnsupportedAccountType
	}

	// 生成6位随机验证码
	code := generateVerificationCode()

	// 发送邮件
	if err := u.emailService.SendVerificationCode(ctx, account, code, purpose); err != nil {
		zlog.CtxErrorf(ctx, "send verification code failed: %v", err)
		return ErrInternalError
	}

	return nil
}

// VerifyCode 校验验证码
func (u *UserServiceImpl) verifyCode(ctx context.Context, account, accountType, code, purpose string) error {
	if account == "" || code == "" || purpose == "" {
		return ErrInvalidParams
	}

	// 从Redis获取验证码
	key := fmt.Sprintf("verification_code:%s:%s", purpose, account)
	storedCode, err := cache.GetRedis(ctx, key)
	if err != nil {
		zlog.CtxErrorf(ctx, "get verification code from redis failed: %v", err)
		return ErrInternalError
	}

	if storedCode == "" {
		zlog.CtxWarnf(ctx, "verification code not found or expired for: %s", account)
		return ErrVerificationCodeIncorrect
	}

	if storedCode != code {
		zlog.CtxWarnf(ctx, "verification code mismatch for: %s", account)
		return ErrVerificationCodeIncorrect
	}

	// 校验成功后删除验证码（一次性使用）
	if err := cache.DelRedis(ctx, key); err != nil {
		zlog.CtxErrorf(ctx, "delete verification code from redis failed: %v", err)
		// 不返回错误，因为验证码已经校验成功
	}

	return nil
}

// generateVerificationCode 生成6位随机验证码
func generateVerificationCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		// crypto/rand 的失败是一个罕见且严重的事件，表明系统的熵源存在问题。
		// 在这种情况下，记录严重错误并 panic 是一个合理的做法。
		panic(fmt.Sprintf("failed to generate cryptographically secure random number for verification code: %v", err))
	}
	return fmt.Sprintf("%06d", n.Int64())
}
