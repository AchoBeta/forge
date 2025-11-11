package oauth

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"

	"forge/infra/configs"
	wechat "forge/infra/oauth/providers"
	"forge/pkg/log/zlog"
)

// InitGoth 初始化 goth OAuth 服务
func InitGoth(oauthConfig configs.OAuthConfig) {
	// 注册 GitHub Provider
	if oauthConfig.GitHubClientID != "" && oauthConfig.GitHubClientSecret != "" {
		goth.UseProviders(
			github.New(
				oauthConfig.GitHubClientID,
				oauthConfig.GitHubClientSecret,
				oauthConfig.GitHubCallbackURL,
			),
		)
		zlog.Infof("GitHub OAuth Provider 初始化成功")
	}

	// 注册微信 Provider
	if oauthConfig.WechatAppID != "" && oauthConfig.WechatAppSecret != "" {
		goth.UseProviders(
			wechat.NewWechat(
				oauthConfig.WechatAppID,
				oauthConfig.WechatAppSecret,
				oauthConfig.WechatCallbackURL,
			),
		)
		zlog.Infof("微信 OAuth Provider 初始化成功")
	}

	// 配置 Session Store
	if oauthConfig.SessionSecret == "" {
		zlog.Warnf("OAuth session_secret 未配置，使用默认值（仅开发环境）")
		oauthConfig.SessionSecret = "default-session-secret-change-in-production-min-32-chars"
	}

	// gorilla/sessions 的 NewCookieStore 可以接受一个或多个密钥
	// 如果只有一个密钥，它会被用于认证和加密
	// 建议使用至少 32 字节的密钥
	secretBytes := []byte(oauthConfig.SessionSecret)

	// 如果密钥长度不足 32 字节，使用哈希扩展
	if len(secretBytes) < 32 {
		zlog.Warnf("OAuth session_secret 长度不足32字节，使用扩展密钥")
		// 重复密钥直到达到 32 字节
		extended := make([]byte, 0, 64)
		for len(extended) < 64 {
			extended = append(extended, secretBytes...)
		}
		secretBytes = extended[:64] // 使用 64 字节的密钥（更安全）
	} else if len(secretBytes) < 64 {
		// 如果密钥长度在 32-64 字节之间，扩展到 64 字节
		extended := make([]byte, 64)
		copy(extended, secretBytes)
		for i := len(secretBytes); i < 64; i++ {
			extended[i] = secretBytes[i%len(secretBytes)]
		}
		secretBytes = extended
	}

	// 使用密钥创建 CookieStore
	// 可以传入多个密钥，第一个用于认证，第二个用于加密（如果提供）
	store := sessions.NewCookieStore(secretBytes)
	store.MaxAge(86400 * 30) // 30天
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false                  // 生产环境改为 true (HTTPS)
	store.Options.SameSite = http.SameSiteLaxMode // Lax 模式允许 GET 请求跨站发送 Cookie（适合 OAuth 回调）

	gothic.Store = store
	zlog.Infof("OAuth Session Store 初始化成功")
}

// GetProvider 获取指定名称的 Provider
func GetProvider(provider string) (goth.Provider, error) {
	return goth.GetProvider(provider)
}
