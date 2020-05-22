package aliyun

import (
	"github.com/JointFaaS/Manager/function"
	"github.com/aliyun/fc-go-sdk"
)

// ListFunction returns metadata of a function
func (m *Manager) ListFunction() ([]*function.Meta, error) {
	out, err := m.fcClient.ListFunctions(fc.NewListFunctionsInput(
		service,
	))

	if err != nil {
		return nil, err
	}
	ret := make([]*function.Meta, len(out.Functions))
	for i, v := range out.Functions {
		ret[i] = &function.Meta{
			FunctionName:         *v.FunctionName,
			MemorySize:           int64(*v.MemorySize),
			Timeout:              int64(*v.Timeout),
			Description:          *v.Description,
			CreatedTime:          *v.CreatedTime,
			CodeChecksum:         *v.CodeChecksum,
			EnvironmentVariables: v.EnvironmentVariables,
			Runtime:              *v.Runtime,
		}
	}
	return ret, nil
}
