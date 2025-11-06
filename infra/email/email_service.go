package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"forge/biz/adapter"
	"forge/infra/configs"
	"forge/pkg/log/zlog"
	templateEmail "forge/template/email"

	"gopkg.in/gomail.v2"
)

type emailServiceImpl struct {
	smtpConfig               configs.SMTPConfig
	verificationCodeTemplate *template.Template
}

var es *emailServiceImpl

func InitEmailService(smtpConfig configs.SMTPConfig) {
	// 使用 embed 包嵌入的模板内容（编译时嵌入到二进制文件中）
	// 模板文件保留在原位置 template/email/verification_code.html
	tmpl, err := template.New("verification_code").Parse(templateEmail.VerificationCodeTemplate)
	if err != nil {
		zlog.Errorf("解析邮件模板失败: %v", err)
		panic(fmt.Sprintf("解析邮件模板失败: %v", err))
	}

	es = &emailServiceImpl{
		smtpConfig:               smtpConfig,
		verificationCodeTemplate: tmpl,
	}
	zlog.Infof("邮件服务初始化成功，模板已嵌入到二进制文件中")
}

func GetEmailService() adapter.EmailService {
	return es
}

// SendVerificationCode 发送邮件 携带验证码
func (e *emailServiceImpl) SendVerificationCode(ctx context.Context, email, code string) error {
	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.smtpConfig.SmtpUser, e.smtpConfig.EncodedName))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "您的验证码")

	// 使用模板渲染HTML邮件内容 唯一变量 Code
	data := map[string]string{
		"Code": code,
	}
	var emailBody bytes.Buffer
	if err := e.verificationCodeTemplate.Execute(&emailBody, data); err != nil {
		zlog.CtxErrorf(ctx, "渲染邮件模板失败: %v", err)
		return fmt.Errorf("渲染邮件模板失败: %w", err)
	}

	m.SetBody("text/html", emailBody.String())

	// 创建邮件发送器
	d := gomail.NewDialer(e.smtpConfig.SmtpHost, e.smtpConfig.SmtpPort, e.smtpConfig.SmtpUser, e.smtpConfig.SmtpPass)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		zlog.CtxErrorf(ctx, "发送邮件失败: %v", err)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	zlog.CtxInfof(ctx, "验证码发送成功，邮箱: %s", email)
	return nil
}
