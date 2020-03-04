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
	DeleteFunction(funcName string) (error)
	GetCodeURI(funcName string) (string, error)
	GetImage(funcName string) (string, error)
	SaveCode(funcName string, file string) (error)
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
	var platformManager PlatformManager
	var err error
	if config.Aliyun.AccessKeyID != "" {
		platformManager, err = aliyun.NewManagerWithConfig(config.Aliyun)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("No available backend")
	}
	sche, err := scheduler.New(platformManager)
	if err != nil {
		return nil, err
	}
	sche.Work()
	m := &Manager{
		scheduler: sche,
		platformManager: platformManager,
	}
	return m, nil
}

