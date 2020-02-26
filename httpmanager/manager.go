package httpmanager

import (
	"errors"

	"github.com/JointFaaS/Manager/env"
	"github.com/JointFaaS/Manager/function"
	"github.com/JointFaaS/Manager/worker"
	"github.com/JointFaaS/Manager/provider/aliyun"
)

// Config includes
type Config struct {
	Aliyun aliyun.Config `yaml:"aliyun"`
}

// PlatformManager is a layer to decouple backend
type PlatformManager interface {
	CreateFunction(funcName string, dir string, e env.Env) (error)
	InvokeFunction(funcName string, args []byte) ([]byte, error)
	ListFunction() ([]*function.Meta, error)
	GetFunction(funcName string) (*function.Meta, error) 
}
// Manager works as an adaptor between JointFaaS and specified cloud
type Manager struct {
	platformManager PlatformManager

	workers map[string]*worker.Worker

	funcToWorker map[string][]*worker.Worker
}

// NewManager builds a manager with given config
func NewManager(config Config) (*Manager, error)  {
	m := &Manager{
		workers: make(map[string]*worker.Worker),
		funcToWorker: make(map[string][]*worker.Worker),
	}
	if config.Aliyun.AccessKeyID != "" {
		aliyunManager, err := aliyun.NewManagerWithConfig(config.Aliyun)
		if err != nil {
			return nil, err
		}
		m.platformManager = aliyunManager
	} else {
		return nil, errors.New("No available backend")
	}
	return m, nil
}

