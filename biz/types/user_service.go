package types

import (
	"context"
	"forge/biz/entity"
)

type IUserService interface {
	Login(ctx context.Context, account, accountType, password string) (*entity.User, string, error) // 返回用户、token、错误

	// Register 基于手机号/邮箱进行注册
	Register(ctx context.Context, req *RegisterParams) (*entity.User, error)

	// ResetPassword 重置密码
	ResetPassword(ctx context.Context, req *ResetPasswordParams) error

	// GetUserByID 根据用户ID获取用户信息（用于JWT鉴权等场景）
	GetUserByID(ctx context.Context, userID string) (*entity.User, error)

	// SendVerificationCode 发送验证码
	// purpose: 使用场景，用于决定账号验证逻辑
	SendVerificationCode(ctx context.Context, account, accountType, purpose string) error

	// UpdateAccount 更新联系方式（绑定/换绑手机号或邮箱）
	UpdateAccount(ctx context.Context, req *UpdateAccountParams) (string, error)

	// UnbindAccount 解绑联系方式（手机号/邮箱）
	UnbindAccount(ctx context.Context, req *UnbindAccountParams) error

	// VerifyCode 验证验证码
	VerifyCode(ctx context.Context, account, accountType, code string) error

	// UpdateAvatar 更新用户头像
	UpdateAvatar(ctx context.Context, userID, avatarURL string) error
}

// 注册参数
type RegisterParams struct {
	UserName    string
	Account     string
	AccountType string // 手机号/邮箱
	Code        string
	Password    string
}

// 重置密码参数
type ResetPasswordParams struct {
	Account         string
	AccountType     string // 手机号/邮箱
	Code            string
	NewPassword     string
	ConfirmPassword string
}

// 更新联系方式参数
type UpdateAccountParams struct {
	Account     string // 新手机号/邮箱
	AccountType string // 手机号/邮箱
	Code        string // 验证码
	Password    string // 密码（如果用户没有密码则必填，如果有密码则可选）
}

// 解绑联系方式参数
type UnbindAccountParams struct {
	Account     string // 需要解绑的手机号/邮箱
	AccountType string // 手机号/邮箱
}

const (
	AccountTypePhone = "phone"
	AccountTypeEmail = "email"
)

// 验证码使用场景
const (
	PurposeRegister      = "register"       // 注册场景
	PurposeResetPassword = "reset_password" // 重置密码场景
	PurposeChangeAccount = "change_account" // 换绑联系方式场景（手机号/邮箱）
)
