package httpmanager

import (
	"errors"

	"github.com/JointFaaS/Manager/env"
	"github.com/JointFaaS/Manager/function"
	"github.com/JointFaaS/Manager/provider/aliyun"
	"github.com/JointFaaS/Manager/worker"
	"github.com/JointFaaS/Manager/scheduler"
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
	GetCodeURI(funcName string) (string, error)
}

type workerSchedule struct {
	funcName string
	resCh chan*worker.Worker
}

// Manager works as an adaptor between JointFaaS and specified cloud
type Manager struct {
	platformManager PlatformManager
	scheduler *scheduler.Scheduler
}

// NewManager builds a manager with given config
func NewManager(config Config) (*Manager, error)  {
	sche, err := scheduler.New()
	if err != nil {
		return nil, err
	}
	m := &Manager{
		scheduler: sche,
	}
	sche.Work()
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

