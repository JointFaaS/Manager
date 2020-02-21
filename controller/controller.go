package controller

import (
	"github.com/JointFaaS/Manager/env"
)

// Config includes
type Config struct {

}

// PlatformManager is a layer to decouple backend
type PlatformManager interface {
	CreateFunction(funcName string, dir string, e env.Env) (error)
	UploadSourceCode(funcName string, dir string, e env.Env) error
	InvokeFunction(funcName string, args *[]byte) ([]byte, error)
}
// Manager works as an adaptor between JointFaaS and specified cloud
type Manager struct {
	platformManager PlatformManager
}

// NewManager builds a manager with given config
func NewManager(config Config) (*Manager, error)  {
	return nil, nil
}

