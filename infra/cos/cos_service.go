package cos

import (
	"fmt"
	"forge/infra/configs"
	"forge/pkg/log/zlog"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type cosServiceImpl struct {
	config *configs.COSConfig
	client *sts.Client
}

var cs *cosServiceImpl

func InitCOSService() {
	cfg := configs.Config().GetCOSConfig()

	client := sts.NewClient(
		cfg.SecretID,
		cfg.SecretKey,
		nil,
	)

	cs = &cosServiceImpl{
		config: &cfg,
		client: client,
	}

	zlog.Infof("COS service initialized successfully, region: %s, bucket: %s", cfg.Region, cfg.Bucket)
}

func GetCOSService() *cosServiceImpl {
	return cs
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
