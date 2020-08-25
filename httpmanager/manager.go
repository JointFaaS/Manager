package httpmanager

import (
	"errors"
	"net/http"

	"github.com/JointFaaS/Manager/env"
	"github.com/JointFaaS/Manager/function"
	"github.com/JointFaaS/Manager/provider/aliyun"
	"github.com/JointFaaS/Manager/provider/aws"
	"github.com/JointFaaS/Manager/scheduler"
	"github.com/JointFaaS/Manager/worker"
)

// HTTPConfig includes the parameter related to http server
type HTTPConfig struct {
	Port string `yaml:"port"`
}

type StorageConfig struct {
	Addr string `yaml:"addr"`
}

// Config includes
type Config struct {
	Aliyun  aliyun.Config `yaml:"aliyun"`
	Aws     aws.Config    `yaml:"aws"`
	Server  HTTPConfig    `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
}

// PlatformManager is a layer to decouple backend
type PlatformManager interface {
	CreateFunction(funcName string, dir string, e env.Env, memory string, timeout string) error
	InvokeFunction(funcName string, args []byte) ([]byte, error)
	ListFunction() ([]*function.Meta, error)
	GetFunction(funcName string) (*function.Meta, error)
	DeleteFunction(funcName string) error
	GetCodeURI(funcName string) (string, error)
	GetImage(funcName string) (string, error)
	SaveCode(funcName string, file string) error
}

type workerSchedule struct {
	funcName string
	resCh    chan *worker.Worker
}

// Manager works as an adaptor between JointFaaS and specified cloud
type Manager struct {
	platformManager PlatformManager
	scheduler       *scheduler.Scheduler

	server *http.ServeMux

	port string

	registrationResponse workerRegistrationResponseBody
}

// NewManager builds a manager with given config
func NewManager(config Config) (*Manager, error) {
	var platformManager PlatformManager
	var rrb workerRegistrationResponseBody
	var err error
	if config.Aliyun.AccessKeyID != "" {
		platformManager, err = aliyun.NewManagerWithConfig(config.Aliyun)
		if err != nil {
			return nil, err
		}
		rrb.AccessKeyID = config.Aliyun.AccessKeyID
		rrb.AccessKeySecret = config.Aliyun.AccessKeySecret
		rrb.JointfaasEnv = "aliyun"
		rrb.Region = config.Aliyun.RegionID
	} else if config.Aws.AccessKeyID != "" {
		platformManager, err = aws.NewManagerWithConfig(config.Aws)
		if err != nil {
			return nil, err
		}
		rrb.AccessKeyID = config.Aws.AccessKeyID
		rrb.AccessKeySecret = config.Aws.AccessKeySecret
		rrb.JointfaasEnv = "aws"
		rrb.Region = config.Aws.RegionID
	} else {
		return nil, errors.New("No available backend")
	}
	rrb.CenterStorage = config.Storage.Addr

	sche, err := scheduler.New(platformManager)
	if err != nil {
		return nil, err
	}
	sche.Work()
	m := &Manager{
		scheduler:            sche,
		platformManager:      platformManager,
		server:               http.NewServeMux(),
		port:                 config.Server.Port,
		registrationResponse: rrb,
	}
	m.setRouter()
	m.setMetrics()
	return m, nil
}
