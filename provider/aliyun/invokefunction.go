package aliyun

import (
	"github.com/aliyun/fc-go-sdk"
)

// InvokeFunction invokes a existed function
func (m *Manager) InvokeFunction(funcName string, args []byte) ([]byte, error) {
	out, err := m.fcClient.InvokeFunction(&fc.InvokeFunctionInput{
		ServiceName:  &service,
		FunctionName: &funcName,
		Payload:      &args,
	})
	if err != nil {
		return nil, err
	}
	return out.Payload, nil
}
