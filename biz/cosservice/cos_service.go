package cosservice

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"forge/biz/adapter"
	"forge/biz/entity"
	"forge/biz/types"
	"forge/infra/configs"
	"forge/pkg/log/zlog"
)

const (
	// MinSTSDuration 腾讯云STS临时密钥最短有效期（秒）
	MinSTSDuration = 900 // 15分钟
	// MaxSTSDuration 腾讯云STS临时密钥最长有效期（秒）
	MaxSTSDuration = 7200 // 2小时
)

// 错误定义
var (
	ErrInvalidParams       = errors.New("参数无效")
	ErrInternalError       = errors.New("内部错误")
	ErrPermissionDenied    = errors.New("权限不足")
	ErrInvalidResourcePath = errors.New("无效的资源路径")
	ErrInvalidDuration     = errors.New("无效的有效期")
)

// COSServiceImpl COS服务实现
type COSServiceImpl struct {
	cosService adapter.COSService
	config     configs.COSConfig
}

func NewCOSServiceImpl(cosService adapter.COSService) *COSServiceImpl {
	cfg := configs.Config().GetCOSConfig()
	return &COSServiceImpl{
		cosService: cosService,
		config:     cfg,
	}
}

// GetOSSCredentials 获取OSS临时凭证
func (s *COSServiceImpl) GetOSSCredentials(ctx context.Context, req *types.GetOSSCredentialsParams) (*types.OSSCredentials, error) {
	// 从JWT token上下文中获取用户信息
	user, ok := entity.GetUser(ctx)
	if !ok {
		zlog.CtxErrorf(ctx, "failed to get user from context")
		return nil, ErrPermissionDenied
	}

	// 参数校验：资源路径
	if req.ResourcePath == "" {
		zlog.CtxErrorf(ctx, "resource path is required")
		return nil, ErrInvalidResourcePath
	}

	// 权限校验：确保资源路径以该用户的路径开头，后续可调整
	expectedPrefix := fmt.Sprintf("user/%s/", user.UserID)
	if !strings.HasPrefix(req.ResourcePath, expectedPrefix) {
		zlog.CtxWarnf(ctx, "resource path does not match user, path: %s, userID: %s", req.ResourcePath, user.UserID)
		return nil, ErrPermissionDenied
	}

	// 参数校验：有效期
	durationSeconds := req.DurationSeconds
	if durationSeconds == 0 {
		durationSeconds = s.config.STSDuration
	}
	if durationSeconds < MinSTSDuration || durationSeconds > MaxSTSDuration {
		zlog.CtxErrorf(ctx, "invalid duration seconds: %d, must be between %d and %d", durationSeconds, MinSTSDuration, MaxSTSDuration)
		return nil, ErrInvalidDuration
	}

	// 调用基础设施层获取临时凭证
	result, err := s.cosService.GetTemporaryCredentials(req.ResourcePath, durationSeconds)
	if err != nil {
		zlog.CtxErrorf(ctx, "failed to get credentials: %v", err)
		return nil, ErrInternalError
	}

	// 解析过期时间
	expiration, err := parseExpiration(result.Expiration)
	if err != nil {
		zlog.CtxErrorf(ctx, "failed to parse expiration: %v", err)
		return nil, ErrInternalError
	}

	// 构建响应
	credentials := &types.OSSCredentials{
		AccessKeyID:     result.Credentials.TmpSecretID,
		SecretAccessKey: result.Credentials.TmpSecretKey,
		SessionToken:    result.Credentials.SessionToken,
		Expiration:      expiration,
		Region:          s.config.Region,
		BucketName:      s.config.Bucket,
		Provider:        "cos",
		BaseURL:         s.config.BaseURL,
		AccountID:       s.config.AppID,
	}

	zlog.CtxInfof(ctx, "OSS credentials retrieved successfully, userID: %s, resource: %s", user.UserID, req.ResourcePath)
	return credentials, nil
}

// parseExpiration 解析过期时间字符串为Unix时间戳
func parseExpiration(expiration string) (int64, error) {
	t, err := time.Parse(time.RFC3339, expiration)
	if err != nil {
		return 0, fmt.Errorf("failed to parse expiration time '%s': %w", expiration, err)
	}
	return t.Unix(), nil
}
