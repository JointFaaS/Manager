package aws

import (
	"github.com/JointFaaS/Manager/env"
	"github.com/JointFaaS/Manager/function"
)

// CreateFunction creates a function on lambda
func (m *Manager) CreateFunction(funcName string, dir string, e env.Env) (error) {
	return nil
}

// GetFunction fetches the metadata of a function
func (m *Manager) GetFunction(funcName string) (*function.Meta, error) {
	return nil, nil	
}

func (m *Manager) InvokeFunction(funcName string, args []byte) ([]byte, error) {
	return nil, nil
}

func (m *Manager) ListFunction() ([]*function.Meta, error) {
	return nil, nil
}

func (m *Manager) DeleteFunction(funcName string) (error) {
	return nil
}
