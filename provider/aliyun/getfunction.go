package aliyun

import (
	"github.com/JointFaaS/Manager/function"
	"github.com/aliyun/fc-go-sdk"
)

// GetFunction returns metadata of a function
func (m *Manager) GetFunction(funcName string) (*function.Meta, error) {
	v, err := m.fcClient.GetFunction(fc.NewGetFunctionInput(
		service,
		funcName,
	))
	if err != nil {
		return nil, err
	}
	return &function.Meta{
		FunctionName:         *v.FunctionName,
		MemorySize:           int64(*v.MemorySize),
		Timeout:              int64(*v.Timeout),
		Description:          *v.Description,
		CreatedTime:          *v.CreatedTime,
		CodeChecksum:         *v.CodeChecksum,
		EnvironmentVariables: v.EnvironmentVariables,
		Runtime:              *v.Runtime,
	}, nil
}
