package aliyun

import (
	"errors"
	"github.com/aliyun/fc-go-sdk"
)

// GetImage returns image of a function
func (m *Manager) GetImage(funcName string) (string, error) {
	funcOut, err := m.fcClient.GetFunction(fc.NewGetFunctionInput(
		service,
		funcName,
	))
	if err != nil {
		return "", err
	}

	runtime := *funcOut.Runtime
	if runtime == "python3" {
		return "registry.cn-shanghai.aliyuncs.com/veia/hcloud-py", nil
	} else if runtime == "java8" {
		return "registry.cn-shanghai.aliyuncs.com/veia/hcloud-java", nil
	} 

	return "", errors.New("Not support env")
}