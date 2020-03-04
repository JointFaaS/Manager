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

	// original code
	userCodeBucket *oss.Bucket
	// aliyun specified code
	aliCodeBucket *oss.Bucket
}

// Config defines the neccessary settings about aliyun.Manager
type Config struct {
	RegionID string `yaml:"regionID"`
	AccessKeyID string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`

	FcEndpoint string `yaml:"fcEndpoint"`
	OssEndpoint string `yaml:"ossEndpoint"`
	UserCodeBucket string `yaml:"userCodeBucket"`
	AliCodeBucket string `yaml:"aliCodeBucket"`
}

// NewManagerWithConfig create a manager with given configuration
func NewManagerWithConfig(config Config) (*Manager, error){
	sdk, err := sdk.NewClientWithAccessKey(config.RegionID, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	ossSdk, err := oss.New(config.OssEndpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	userBucket, err := ossSdk.Bucket(config.UserCodeBucket)
	if err != nil {
		return nil, err
	}
	aliBucket, err := ossSdk.Bucket(config.AliCodeBucket)
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
		userCodeBucket: userBucket,
		aliCodeBucket: aliBucket,
		fcClient: fcSdk,
	}
	return manager, nil
}