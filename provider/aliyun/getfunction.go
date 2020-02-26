package aliyun

import (
	"github.com/JointFaaS/Manager/function"
	"github.com/aliyun/fc-go-sdk"
)

// GetFunction returns metadata of a function
func (m *Manager) GetFunction(funcName string) (*function.Meta, error) {
	funcOut, err := m.fcClient.GetFunction(fc.NewGetFunctionInput(
		service,
		funcName,
	))
	if err != nil {
		return nil, err
	}
	return &function.Meta{
		FunctionName: funcName,
		Description: *funcOut.Description,
		Runtime: *funcOut.Runtime,
		// TODO
	}, nil
}