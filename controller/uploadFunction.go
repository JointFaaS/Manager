package controller

import (
	"github.com/JointFaaS/Manager/env"
)

// UploadFunctionInput metadata
type UploadFunctionInput struct {

}

// UploadFunction creates a new function
func (m* Manager) UploadFunction(funcName string, dir string, env env.Env) (error) {
	err := m.platformManager.CreateFunction(funcName, dir, env)
	if err != nil {
		return err
	}
	return nil
}

