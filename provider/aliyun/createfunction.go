package aliyun

import (
	"errors"
	"path"
	"strconv"

	"github.com/JointFaaS/Manager/env"
	"github.com/aliyun/fc-go-sdk"
)

var service = "jointfaas"

// CreateFunction :
// sourceURL can be created by UploadSourceCode
func (m *Manager) CreateFunction(funcName string, dir string, e env.Env, memoryS string, timeoutS string) (error) {
	var err error
	memorySizeI, err := strconv.Atoi(memoryS)
	memorySize := int32(memorySizeI)
	if err != nil {
		return err
	}
	timeoutI, err := strconv.Atoi(timeoutS)
	timeout := int32(timeoutI)
	if err != nil {
		return err
	}
	if e == env.PYTHON3 {
		err = m.createPython3Function(path.Join(dir, "code"))
	}else if e == env.JAVA8 {
		err = m.createJava8Function(path.Join(dir, "code"))
	} else {
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

	err = m.aliCodeBucket.PutObjectFromFile(funcName, aliyunZip)
	if err != nil {
		return err
	}
	runtime, handler, initializer := envToAliyunEnv(e)
	_, err = m.fcClient.CreateFunction(&fc.CreateFunctionInput{
		ServiceName: &service,
		FunctionCreateObject: fc.FunctionCreateObject{
			FunctionName: &funcName,
			Runtime: &runtime,
			Initializer: &initializer,
			Handler: &handler,
			Code: fc.NewCode().WithOSSBucketName(m.aliCodeBucket.BucketName).WithOSSObjectName(funcName),
			MemorySize: &memorySize,
			Timeout: &timeout,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func envToAliyunEnv(e env.Env) (string, string, string) {
	if e == env.PYTHON3 {
		return "python3", "jointfaas.handler", ""
	} else if e == env.JAVA8 {
		return "java8", "jointfaas.AliIndex::handleRequest", "jointfaas.AliIndex::initialize"
	}
	return "", "", ""
}
