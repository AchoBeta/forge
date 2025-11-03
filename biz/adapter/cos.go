package adapter

import (
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

// COSService 定义COS服务的业务接口
type COSService interface {
	GetTemporaryCredentials(resourcePath string, durationSeconds int64) (*sts.CredentialResult, error)
}
