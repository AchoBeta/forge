package email

import (
	"context"
	"fmt"

	"forge/biz/adapter"
	"forge/infra/configs"
	"forge/pkg/log/zlog"

	"gopkg.in/gomail.v2"
)

type emailServiceImpl struct {
	smtpConfig configs.IConfig
}

var es *emailServiceImpl

func InitEmailService(config configs.IConfig) {
	es = &emailServiceImpl{
		smtpConfig: config,
	}
}

func GetEmailService() adapter.EmailService {
	return es
}

// SendVerificationCode 发送邮件 携带验证码
func (e *emailServiceImpl) SendVerificationCode(ctx context.Context, email, code, purpose string) error {
	smtpCfg := e.smtpConfig.GetSMTPConfig()

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(smtpCfg.SmtpUser, smtpCfg.EncodedName))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "您的验证码")

	// 构建HTML邮件内容
	emailBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<style>
			body {
				font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
				line-height: 1.6;
				color: #333;
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
			}
			.container {
				border: 1px solid #eaeaea;
				border-radius: 5px;
				padding: 20px;
				background-color: #ffffff;
			}
			h2 {
				color: #333;
				margin-top: 0;
			}
			.code-box {
				font-size: 24px;
				font-weight: bold;
				letter-spacing: 2px;
				color: #1890ff;
				margin: 20px 0;
				padding: 15px;
				background-color: #f5f5f5;
				border-radius: 4px;
				display: inline-block;
				text-align: center;
			}
			.footer {
				font-size: 14px;
				color: #999;
				margin-top: 20px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>邮箱验证码</h2>
			<p>您的验证码是：</p>
			<div class="code-box">%s</div>
			<p class="footer">此验证码10分钟内有效，请勿泄露给他人。</p>
		</div>
	</body>
	</html>
	`, code)

	m.SetBody("text/html", emailBody)

	// 创建邮件发送器
	d := gomail.NewDialer(smtpCfg.SmtpHost, smtpCfg.SmtpPort, smtpCfg.SmtpUser, smtpCfg.SmtpPass)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		zlog.CtxErrorf(ctx, "发送邮件失败: %v", err)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	zlog.CtxInfof(ctx, "验证码发送成功，邮箱: %s, 用途: %s", email, purpose)
	return nil
}
