package aws

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/JointFaaS/Manager/env"
	"github.com/JointFaaS/Manager/function"
)

// CreateFunction creates a function on lambda
func (m *Manager) CreateFunction(funcName string, dir string, e env.Env) (error) {
	if e == env.PYTHON3 {
		m.injectPython3Handler(path.Join(dir, "code"))
	} else if e == env.JAVA8 {

	} else {
		return errors.New("Not support env")
	}

	awsZip := path.Join(dir, "aws.zip")
	err := compressDir(path.Join(dir, "code"), awsZip)
	if err != nil {
		return err
	}
	body, err := os.Open(awsZip)
	if err != nil {
		return err
	}

	_, err = m.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: &m.awsCodeBucket,
		Key: &funcName,
		Body: body,
	})
	if err != nil {
		return err
	}
	
	runtime, handler := envToAWSEnv(e)

	// arn:${Partition}:lambda:${Region}:${Account}:function:${FunctionName}
	role := fmt.Sprintf("arn:%s:lambda:%s:%s:function:%s", m.lambdaClient.PartitionID, *m.lambdaClient.Config.Region, m.account, funcName)
	_, err = m.lambdaClient.CreateFunction(&lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			S3Bucket: &m.awsCodeBucket,
			S3Key: &funcName,
		},
		FunctionName: &funcName,
		Runtime: &runtime,
		Handler: &handler,
		Role: &role,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetFunction fetches the metadata of a function
func (m *Manager) GetFunction(funcName string) (*function.Meta, error) {
	output, err := m.lambdaClient.GetFunction(&lambda.GetFunctionInput{
		FunctionName: &funcName,
	})
	if err != nil {
		return nil, err
	}
	f := output.Configuration
	return &function.Meta{
		FunctionName: funcName,
		Runtime: *f.Runtime,	// TODO: translate
		Handler: *f.Handler,
		MemorySize: *f.MemorySize,
		CreatedTime: *f.LastModified,
	}, nil	
}

// InvokeFunction invokes a function by name
func (m *Manager) InvokeFunction(funcName string, args []byte) ([]byte, error) {
	output, err := m.lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName: &funcName,
		Payload: args,
	})
	if err != nil {
		return nil, err
	}
	return output.Payload, nil
}

// ListFunction displays the functions
func (m *Manager) ListFunction() ([]*function.Meta, error) {
	functions := make([]*function.Meta, 0)
	for {
		var marker *string = nil
		output, err := m.lambdaClient.ListFunctions(&lambda.ListFunctionsInput{Marker: marker})
		if err != nil {
			return nil, err
		}
		for _, f := range output.Functions {
			functions = append(functions, &function.Meta{
				FunctionName: *f.FunctionName,
				Runtime: *f.Runtime,	// TODO: translate
				Handler: *f.Handler,
				MemorySize: *f.MemorySize,
				CreatedTime: *f.LastModified,
			})
		}
		marker = output.NextMarker
		if marker == nil {
			break
		}
	}
	return functions, nil
}

// DeleteFunction deletes a function
func (m *Manager) DeleteFunction(funcName string) (error) {
	_, err := m.lambdaClient.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: &funcName,
	})
	return err
}

func envToAWSEnv(e env.Env) (string, string) {
	if e == env.PYTHON3 {
		return "python3", "jointfaas.handler"
	} else if e == env.JAVA8 {
		return "java8", "jointfaas.AliIndex::handleRequest"
	}
	return "", ""
}