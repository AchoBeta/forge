package cos

import (
	"fmt"
	"forge/biz/adapter"
	"forge/infra/configs"
	"forge/pkg/log/zlog"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type cosServiceImpl struct {
	config configs.COSConfig
	client *sts.Client
}

// NewCOSService 创建COS服务实例（依赖注入模式）
// 通过构造函数接收配置，返回接口类型，便于测试和依赖注入
func NewCOSService(cfg configs.COSConfig) adapter.COSService {
	client := sts.NewClient(
		cfg.SecretID,
		cfg.SecretKey,
		nil,
	)

	service := &cosServiceImpl{
		config: cfg,
		client: client,
	}

	zlog.Infof("COS service created successfully, region: %s, bucket: %s", cfg.Region, cfg.Bucket)
	return service
}

// GetTemporaryCredentials 获取COS临时凭证
func (c *cosServiceImpl) GetTemporaryCredentials(resourcePath string, durationSeconds int64) (*sts.CredentialResult, error) {
	// 构建资源ARN
	resourceArn := fmt.Sprintf(
		"qcs::cos:%s:uid/%s:%s-%s/%s",
		c.config.Region,
		c.config.AppID,
		c.config.Bucket,
		c.config.AppID,
		resourcePath,
	)

	// 配置STS策略
	opt := &sts.CredentialOptions{
		DurationSeconds: durationSeconds,
		Region:          c.config.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						// 简单上传操作
						"name/cos:PostObject",
						"name/cos:PutObject",
						"name/cos:GetObject",
						"name/cos:HeadObject",
						// 分片上传操作
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						resourceArn,
					},
				},
			},
		},
	}

	// 请求临时凭证
	result, err := c.client.GetCredential(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to get COS STS credentials: %w", err)
	}

	return result, nil
}
