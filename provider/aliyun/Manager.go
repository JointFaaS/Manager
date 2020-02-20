package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
)

// Manager exposes the interfaces about resource CRUD on aliyun
type Manager struct {
	sdkClient *sdk.Client
	ossClient *oss.Client

	fcClient *fc.Client

	codeBucket *oss.Bucket
}

// Config defines the neccessary settings about aliyun.Manager
type Config struct {
	regionID string
	accessKeyID string
	accessKeySecret string

	bucket string
}

// NewManagerWithConfig create a manager with gived configuration
func NewManagerWithConfig(config Config) (*Manager, error){
	sdk, err := sdk.NewClientWithAccessKey(config.regionID, config.accessKeyID, config.accessKeySecret)
	if err != nil {
		return nil, err
	}
	ossSdk, err := oss.New(config.regionID, config.accessKeyID, config.accessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := ossSdk.Bucket(config.bucket)
	if err != nil {
		return nil, err
	}
	fcSdk, err := fc.NewClient(config.regionID, "", config.accessKeyID, config.accessKeySecret)
	if err != nil {
		return nil, err
	}

	fcSdk.CreateService(fc.NewCreateServiceInput().WithServiceName(service))
	manager := &Manager{
		sdkClient: sdk,
		ossClient: ossSdk,
		codeBucket: bucket,
		fcClient: fcSdk,
	}
	return manager, nil
}