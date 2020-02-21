package aliyun

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"io"
	"os"
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
		err = m.createPython3Function(dir)
	}else {
		return errors.New("Not support Env")
	}
	if err != nil {
		return err
	}

	d, _ := os.Create(path.Join(dir, "code.zip"))
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	fileInfos, _ := ioutil.ReadDir(dir)
	for _, fi := range fileInfos {
		f, err := os.Open(path.Join(dir, fi.Name()))
		if err != nil {
			return err
		}
		err = compress(f, "", w)
		if err != nil {
			return err
		}
	}
	d.Seek(0, 0)
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

// InvokeFunction invokes a existed function
func (m *Manager) InvokeFunction(funcName string, args *[]byte) ([]byte, error) {
	_, err := m.fcClient.InvokeFunction(&fc.InvokeFunctionInput{
		ServiceName: &service,
		FunctionName: &funcName,
		Payload: args,
	})
	return nil, err
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func envToAliyunRuntime(e env.Env) string {
	if e == env.PYTHON3 {
		return "python3"
	} else {
		return ""
	}
}