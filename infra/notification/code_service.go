package notification

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"forge/biz/adapter"
	"forge/infra/configs"
	"forge/pkg/log/zlog"
	templateEmail "forge/template/email"

	"gopkg.in/gomail.v2"
)

type codeServiceImpl struct {
	smtpConfig               configs.SMTPConfig
	smsConfig                configs.SMSConfig
	verificationCodeTemplate *template.Template
	httpClient               *http.Client
}

var cs *codeServiceImpl

// InitCodeService 初始化验证码服务，需在程序启动时调用
func InitCodeService(smtpConfig configs.SMTPConfig, smsConfig configs.SMSConfig) {
	tmpl, err := template.New("verification_code").Parse(templateEmail.VerificationCodeTemplate)
	if err != nil {
		zlog.Errorf("解析验证码邮件模板失败: %v", err)
		panic(fmt.Sprintf("解析验证码邮件模板失败: %v", err))
	}

	cs = &codeServiceImpl{
		smtpConfig:               smtpConfig,
		smsConfig:                smsConfig,
		verificationCodeTemplate: tmpl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	zlog.Infof("验证码服务初始化成功，已配置邮件与短信通道")
}

// GetCodeService 获取验证码服务实例
func GetCodeService() adapter.CodeService {
	return cs
}

// SendEmailCode 发送邮件验证码
func (c *codeServiceImpl) SendEmailCode(ctx context.Context, email, code string) error {
	if c == nil {
		return fmt.Errorf("code service not initialized")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(c.smtpConfig.SmtpUser, c.smtpConfig.EncodedName))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "您的验证码")

	data := map[string]string{
		"Code": code,
	}
	var emailBody bytes.Buffer
	if err := c.verificationCodeTemplate.Execute(&emailBody, data); err != nil {
		zlog.CtxErrorf(ctx, "渲染验证码邮件模板失败: %v", err)
		return fmt.Errorf("渲染验证码邮件模板失败: %w", err)
	}

	m.SetBody("text/html", emailBody.String())

	d := gomail.NewDialer(c.smtpConfig.SmtpHost, c.smtpConfig.SmtpPort, c.smtpConfig.SmtpUser, c.smtpConfig.SmtpPass)

	if err := d.DialAndSend(m); err != nil {
		zlog.CtxErrorf(ctx, "发送验证码邮件失败: %v", err)
		return fmt.Errorf("发送验证码邮件失败: %w", err)
	}

	zlog.CtxInfof(ctx, "验证码邮件发送成功，邮箱: %s", email)
	return nil
}

// SendSMSCode 发送短信验证码
func (c *codeServiceImpl) SendSMSCode(ctx context.Context, phone, code string) error {
	if c == nil {
		return fmt.Errorf("code service not initialized")
	}

	if c.smsConfig.TemplateID == "" {
		return fmt.Errorf("sms template id not configured")
	}

	endpoint := c.smsConfig.Endpoint
	if endpoint == "" {
		return fmt.Errorf("sms endpoint not configured")
	}

	smsURL := fmt.Sprintf(endpoint, c.smsConfig.TemplateID, url.QueryEscape(code), url.QueryEscape(phone))
	resp, err := c.httpClient.Get(smsURL)
	if err != nil {
		zlog.CtxErrorf(ctx, "请求短信服务失败: %v", err)
		return fmt.Errorf("request sms service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		zlog.CtxErrorf(ctx, "短信服务返回状态码 %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		return fmt.Errorf("sms service returned status %d", resp.StatusCode)
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	zlog.CtxInfof(ctx, "短信验证码发送成功，手机号: %s", phone)
	return nil
}
