package adapter

import "context"

// EmailService 邮箱服务接口
// 用于发送验证码邮件
type EmailService interface {
	// SendVerificationCode 发送验证码到指定邮箱
	// email: 目标邮箱地址
	// code: 验证码
	SendVerificationCode(ctx context.Context, email, code string) error
}
