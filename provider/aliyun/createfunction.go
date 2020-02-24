package aliyun

import (
	"errors"
	"path"

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
		err = m.createPython3Function(path.Join(dir, "code"))
	}else {
		return errors.New("Not support Env")
	}
	if err != nil {
		return err
	}

	aliyunZip := path.Join(dir, "aliyun.zip")
	err = compressDir(path.Join(dir, "code"), aliyunZip)
	if err != nil {
		return err
	}

	err = m.codeBucket.PutObjectFromFile(funcName, aliyunZip)
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