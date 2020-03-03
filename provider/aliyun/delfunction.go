package aliyun

import (
	"github.com/aliyun/fc-go-sdk"
)

// DeleteFunction deletes a given function
func (m *Manager) DeleteFunction(funcName string) error {
	_, err := m.fcClient.DeleteFunction(fc.NewDeleteFunctionInput(
		service,
		funcName,
	))

	return err
}
