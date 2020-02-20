package aliyun

import (
	"bytes"
	"os"
	"github.com/JointFaaS/Manager/env"
	"github.com/aliyun/fc-go-sdk"
)

var handler = "jointfaas.handler"
var service = "jointfaas"

// CreateFunction :
// sourceURL can be created by UploadSourceCode
func (m *Manager) CreateFunction(funcName string, runtime string) (error) {
	_, err := m.fcClient.CreateFunction(&fc.CreateFunctionInput{
		ServiceName: &service,
		FunctionCreateObject: fc.FunctionCreateObject{
			FunctionName: &funcName,
			Runtime: &runtime,
			Handler: &handler,
			Code: fc.NewCode().WithOSSBucketName(service).WithOSSObjectName(funcName),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// InvokeFunction invokes a existed function
func (m *Manager) InvokeFunction(funcName string, args *[]byte) (*fc.InvokeFunctionOutput, error) {
	invokeOutput, err := m.fcClient.InvokeFunction(&fc.InvokeFunctionInput{
		ServiceName: &service,
		FunctionName: &funcName,
		Payload: args,
	})
	return invokeOutput, err
}

// UploadSourceCode receives zip file and upload into oss
func (m *Manager) UploadSourceCode(funcName string, zipFile *os.File) (error) {
	err := m.codeBucket.PutObject(funcName, zipFile)
	if err != nil {
		return err
	}
	return nil
}

// ConvertFormat converts the standard source format in JointFaaS to specified platform format
func (m *Manager) ConvertFormat(zipFile *os.File, handler string, e env.Env) (*os.File, error)  {
	return nil, nil
}

