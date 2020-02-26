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
	RegionID string `yaml:"regionID"`
	AccessKeyID string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`

	FcEndpoint string `yaml:"fcEndpoint"`
	OssEndpoint string `yaml:"ossEndpoint"`
	Bucket string `yaml:"bucket"`
}

// NewManagerWithConfig create a manager with gived configuration
func NewManagerWithConfig(config Config) (*Manager, error){
	sdk, err := sdk.NewClientWithAccessKey(config.RegionID, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	ossSdk, err := oss.New(config.OssEndpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := ossSdk.Bucket(config.Bucket)
	if err != nil {
		return nil, err
	}
	fcSdk, err := fc.NewClient(config.FcEndpoint, "2016-08-15", config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	
	manager := &Manager{
		sdkClient: sdk,
		ossClient: ossSdk,
		codeBucket: bucket,
		fcClient: fcSdk,
	}
	return manager, nil
}