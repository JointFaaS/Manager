package aliyun

import (
	"errors"

	"github.com/JointFaaS/Manager/env"
	"github.com/aliyun/fc-go-sdk"
)

var handler = "jointfaas.handler"
var service = "jointfaas"


// CreateFunction :
// sourceURL can be created by UploadSourceCode
func (m *Manager) CreateFunction(funcName string, dir string, e env.Env) (error) {
	var err error
	if e == env.PYTHON3 {
		err = m.createPython3Function(dir)
	}else {
		return errors.New("Not support Env")
	}
	if err != nil {
		return err
	}

	d, err := compressDir(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	err = m.codeBucket.PutObject(funcName, d)
	if err != nil {
		return err
	}
	runtime := envToAliyunRuntime(e)
	_, err = m.fcClient.CreateFunction(&fc.CreateFunctionInput{
		ServiceName: &service,
		FunctionCreateObject: fc.FunctionCreateObject{
			FunctionName: &funcName,
			Runtime: &runtime,
			Handler: &handler,
			Code: fc.NewCode().WithOSSBucketName(m.codeBucket.BucketName).WithOSSObjectName(funcName),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func envToAliyunRuntime(e env.Env) string {
	if e == env.PYTHON3 {
		return "python3"
	}
	return ""
}