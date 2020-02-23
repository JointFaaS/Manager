package httpmanager

import (
	"errors"

	"github.com/JointFaaS/Manager/env"
	"github.com/JointFaaS/Manager/provider/aliyun"
)

// Config includes
type Config struct {
	Aliyun aliyun.Config `yaml:"aliyun"`
}

// PlatformManager is a layer to decouple backend
type PlatformManager interface {
	CreateFunction(funcName string, dir string, e env.Env) (error)
	InvokeFunction(funcName string, args *[]byte) ([]byte, error)
}
// Manager works as an adaptor between JointFaaS and specified cloud
type Manager struct {
	platformManager PlatformManager
}

// NewManager builds a manager with given config
func NewManager(config Config) (*Manager, error)  {
	m := Manager{}
	if config.Aliyun.AccessKeyID != "" {
		aliyunManager, err := aliyun.NewManagerWithConfig(config.Aliyun)
		if err != nil {
			return nil, err
		}
		m.platformManager = aliyunManager
	} else {
		return nil, errors.New("No available backend")
	}
	return nil, nil
}

