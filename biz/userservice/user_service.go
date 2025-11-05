package userservice

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"path"
	"strings"
	"time"

	"forge/biz/adapter"
	"forge/biz/entity"
	"forge/biz/repo"
	"forge/biz/types"
	"forge/constant"
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
	if err := u.VerifyCode(ctx, req.Account, req.AccountType, req.Code); err != nil {
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

	// 校验验证码 code（短信/邮箱）
	if err := u.VerifyCode(ctx, req.Account, req.AccountType, req.Code); err != nil {
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
func (u *UserServiceImpl) SendVerificationCode(ctx context.Context, account, accountType string) error {
	// 参数校验
	if account == "" || accountType == "" {
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

	// 先将验证码存储到 Redis，并设置过期时间
	key := fmt.Sprintf(constant.REDIS_VERIFICATION_CODE_KEY, account)
	// TODO: 建议将过期时间（10分钟）配置化
	expiration := 10 * time.Minute
	if err := cache.SetRedis(ctx, key, code, expiration); err != nil {
		zlog.CtxErrorf(ctx, "存储验证码到Redis失败: %v", err)
		return ErrInternalError
	}
	// 发送邮件
	if err := u.emailService.SendVerificationCode(ctx, account, code); err != nil {
		zlog.CtxErrorf(ctx, "send verification code failed: %v", err)
		// 邮件发送失败，尝试从Redis中删除已存储的验证码，以保持一致性
		if delErr := cache.DelRedis(ctx, key); delErr != nil {
			zlog.CtxErrorf(ctx, "删除Redis中未发送成功的验证码失败: %v", delErr)
		}
		return ErrInternalError
	}

	return nil
}

// VerifyCode 校验验证码
func (u *UserServiceImpl) VerifyCode(ctx context.Context, account, accountType, code string) error {
	if account == "" || code == "" {
		return ErrInvalidParams
	}

	// 从Redis获取验证码
	key := fmt.Sprintf(constant.REDIS_VERIFICATION_CODE_KEY, account)
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

// UpdateAvatar 更新用户头像
func (u *UserServiceImpl) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
	// 参数校验
	if userID == "" || avatarURL == "" {
		zlog.CtxErrorf(ctx, "invalid params for update avatar: userID or avatarURL is empty")
		return ErrInvalidParams
	}

	// URL验证
	if err := validateAvatarURL(ctx, avatarURL); err != nil {
		zlog.CtxErrorf(ctx, "avatar URL validation failed: %v", err)
		// 包装错误以保留详细验证信息，同时仍可用 errors.Is 检查错误类型
		return fmt.Errorf("%w: %v", ErrInvalidParams, err) // 保留详细错误
	}

	// 检查用户是否存在（GetUserByID 包含状态检查）
	_, err := u.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// 更新头像
	updateInfo := &repo.UserUpdateInfo{
		UserID: userID,
		Avatar: &avatarURL,
	}
	if err := u.userRepo.UpdateUser(ctx, updateInfo); err != nil {
		zlog.CtxErrorf(ctx, "update avatar failed: %v", err)
		return ErrInternalError
	}

	zlog.CtxInfof(ctx, "update avatar successfully for user: %s", userID)
	return nil
}

// validateAvatarURL URL验证函数
// 注意：移除了路径格式强制检查（原 /user/{userID}/avatar/），允许使用外部服务
// 如果需要对自有存储路径进行限制，应该在存储访问层（COS IAM策略）实现
func validateAvatarURL(ctx context.Context, avatarURL string) error {
	// 1. URL长度限制（防止过长的URL）
	const maxURLLength = 2048 // RFC 7230 建议的最大URL长度
	if len(avatarURL) > maxURLLength {
		return fmt.Errorf("avatar URL too long: exceeds %d characters", maxURLLength)
	}

	// 2. 使用标准库解析URL
	parsedURL, err := url.Parse(avatarURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// 3. 验证协议（只允许http和https）
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("invalid URL scheme: only http and https are allowed, got %s", scheme)
	}

	// 4. 验证Host不为空
	if parsedURL.Host == "" {
		return fmt.Errorf("invalid URL: host is required")
	}

	// 5. 验证Host格式（不能包含危险字符）
	// 注意：移除了对 ".." 的检查，因为主机名中的 ".." 不是安全问题（路径遍历发生在路径部分）
	// 虽然 url.Parse 通常会处理 "//"，但保留检查以防格式错误
	if strings.Contains(parsedURL.Host, "//") {
		return fmt.Errorf("invalid URL: host contains invalid characters")
	}

	// 6. SSRF 防护：禁止访问内网/私有IP地址
	// 使用 Hostname() 方法提取主机名，自动处理端口和 IPv6 方括号
	host := parsedURL.Hostname()

	// 解析 IP 地址
	ip := net.ParseIP(host)
	if ip != nil {
		// 如果是 IP 地址，检查是否为私有/保留地址
		if isPrivateIP(ip) {
			return fmt.Errorf("invalid URL: private/internal IP addresses are not allowed for security reasons")
		}
	} else {
		// 如果是域名，解析为 IP 并检查
		ips, err := net.LookupIP(host)
		if err != nil {
			// 域名解析失败，拒绝URL（可能是恶意域名或网络问题）
			zlog.CtxErrorf(ctx, "failed to resolve host %s: %v", host, err)
			return fmt.Errorf("invalid URL: failed to resolve host %s", host)
		}

		// 检查所有解析出的 IP 地址
		if len(ips) == 0 {
			return fmt.Errorf("invalid URL: host %s resolves to no IP addresses", host)
		}

		for _, resolvedIP := range ips {
			if isPrivateIP(resolvedIP) {
				return fmt.Errorf("invalid URL: host %s resolves to private/internal IP address", host)
			}
		}
	}

	// 7. 验证路径中不能包含危险字符（防止路径遍历攻击）
	if strings.Contains(parsedURL.Path, "..") || strings.Contains(parsedURL.Path, "//") {
		return fmt.Errorf("invalid URL path: contains dangerous characters")
	}

	// 8. 允许查询参数（外部服务如 Gravatar、CDN 需要查询参数）
	// 但禁止锚点（Fragment），因为锚点不会发送到服务器
	if parsedURL.Fragment != "" {
		return fmt.Errorf("invalid URL: fragment is not allowed")
	}

	// 9. 验证URL路径或查询参数中是否包含图片格式标识
	// 支持多种常见格式：
	// - 直接路径：https://example.com/avatar.jpg
	// - 查询参数：https://gravatar.com/avatar/xxx?s=200&d=identicon
	// - 路径+查询：https://cdn.example.com/user123.jpg?width=200

	// 从路径中提取可能的文件名
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	var fileName string
	if len(pathParts) > 0 {
		fileName = pathParts[len(pathParts)-1]
	}

	// 检查路径中的文件扩展名
	hasValidExtension := false
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg"}
	// 允许的图片格式（不带点，用于查询参数）
	validImageFormats := []string{"png", "jpg", "jpeg", "gif", "webp", "bmp", "svg"}

	if fileName != "" {
		// 使用 path.Ext 提取真正的文件扩展名，避免被恶意文件名绕过（如 avatar.jpg.exe）
		fileExt := strings.ToLower(path.Ext(fileName))
		for _, ext := range allowedExtensions {
			if fileExt == ext {
				hasValidExtension = true
				break
			}
		}
	}

	// 如果路径中没有有效的扩展名，检查查询参数中是否有图片相关的标识
	// 例如：?format=png, ?type=image 等（某些服务使用查询参数指定格式）
	if !hasValidExtension && parsedURL.RawQuery != "" {
		// 解析查询参数，避免误判（如 ?some_other_param=format=png 不应该被识别）
		// url.Values.Get() 只返回指定键的值，不会因为参数值中包含字符串而误判
		query := parsedURL.Query()

		// 检查 format 参数（如 ?format=png）
		if format := strings.ToLower(query.Get("format")); format != "" {
			for _, validFormat := range validImageFormats {
				if format == validFormat {
					hasValidExtension = true
					break
				}
			}
		}

		// 检查 type 参数（如 ?type=image）
		if !hasValidExtension && strings.ToLower(query.Get("type")) == "image" {
			hasValidExtension = true
		}

		// 检查 mime 参数（如 ?mime=image/png）
		if !hasValidExtension && strings.Contains(strings.ToLower(query.Get("mime")), "image") {
			hasValidExtension = true
		}

		// 检查 ext 参数（如 ?ext=png）
		if !hasValidExtension {
			if ext := strings.ToLower(query.Get("ext")); ext != "" {
				for _, validExt := range validImageFormats {
					if ext == validExt {
						hasValidExtension = true
						break
					}
				}
			}
		}
	}

	// 如果既没有路径扩展名，也没有查询参数标识，允许通过但记录警告
	// 因为某些服务可能通过 Content-Type 响应头来标识图片，而不是URL
	if !hasValidExtension {
		zlog.CtxWarnf(ctx, "avatar URL does not contain explicit image format identifier: %s", avatarURL)
		// 不返回错误，允许通过，因为某些合法的图片URL可能没有扩展名
	}

	// 10. 如果路径中有文件名，验证文件名格式
	if fileName != "" {
		// 验证文件名长度（防止过长的文件名）
		const maxFileNameLength = 255
		if len(fileName) > maxFileNameLength {
			return fmt.Errorf("invalid filename: too long, exceeds %d characters", maxFileNameLength)
		}

		// 验证文件名不能包含明显的危险字符
		// 注意：这里不禁止 : 和 ?，因为它们可能在合法的URL中出现
		dangerousChars := []string{"<", ">", "|", "\"", "*", "\\", "\x00"}
		for _, char := range dangerousChars {
			if strings.Contains(fileName, char) {
				return fmt.Errorf("invalid filename: contains dangerous character '%s'", char)
			}
		}
	}

	return nil
}

// isPrivateIP 检查 IP 地址是否为私有/保留地址（用于 SSRF 防护）
func isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}

	// 使用标准库函数检查常见的私有/保留地址范围（同时支持 IPv4 和 IPv6）
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsPrivate() || ip.IsMulticast() {
		return true
	}

	// 标准库的 IsUnspecified() 只检查单个地址（0.0.0.0 或 ::），但对于 SSRF 防护，
	// 我们应该拒绝整个 0.0.0.0/8 范围（0.0.0.0 到 0.255.255.255）
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 0
	}

	// 对于 IPv6，IsUnspecified() 已足够检查未指定地址（::）
	return ip.IsUnspecified()
}
