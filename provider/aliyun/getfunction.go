package aliyun

import (
	"github.com/JointFaaS/Manager/function"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
)

// GetFunction returns metadata of a function
func (m *Manager) GetFunction(funcName string) (*function.Meta, error) {
	funcOut, err := m.fcClient.GetFunction(fc.NewGetFunctionInput(
		service,
		funcName,
	))
	codeOut, err := m.codeBucket.SignURL(funcName, oss.HTTPGet, 99999999)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &function.Meta{
		FunctionName: funcName,
		Description: *funcOut.Description,
		Runtime: *funcOut.Runtime,
		CodeURI: codeOut,
		// TODO
	}, nil
}