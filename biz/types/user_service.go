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
	SendVerificationCode(ctx context.Context, account, accountType string) error
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

const (
	AccountTypePhone = "phone"
	AccountTypeEmail = "email"
)
