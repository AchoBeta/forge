package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"

	"forge/biz/adapter"
	"forge/infra/configs"
	"forge/pkg/log/zlog"
	"forge/util"

	"gopkg.in/gomail.v2"
)

type emailServiceImpl struct {
	smtpConfig               configs.IConfig
	verificationCodeTemplate *template.Template
}

var es *emailServiceImpl

func InitEmailService(config configs.IConfig) {
	// 加载验证码邮件模板
	templatePath := filepath.Join(util.GetRootPath(""), "template", "email", "verification_code.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		zlog.Errorf("加载邮件模板失败: %v, 模板路径: %s", err, templatePath)
		panic(fmt.Sprintf("加载邮件模板失败: %v", err))
	}

	es = &emailServiceImpl{
		smtpConfig:               config,
		verificationCodeTemplate: tmpl,
	}
	zlog.Infof("邮件服务初始化成功，模板路径: %s", templatePath)
}

func GetEmailService() adapter.EmailService {
	return es
}

// SendVerificationCode 发送邮件 携带验证码
func (e *emailServiceImpl) SendVerificationCode(ctx context.Context, email, code string) error {
	smtpCfg := e.smtpConfig.GetSMTPConfig()

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(smtpCfg.SmtpUser, smtpCfg.EncodedName))
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
	d := gomail.NewDialer(smtpCfg.SmtpHost, smtpCfg.SmtpPort, smtpCfg.SmtpUser, smtpCfg.SmtpPass)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		zlog.CtxErrorf(ctx, "发送邮件失败: %v", err)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	zlog.CtxInfof(ctx, "验证码发送成功，邮箱: %s", email)
	return nil
}
