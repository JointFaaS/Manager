package openstack

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/JointFaaS/Manager/env"
	"github.com/JointFaaS/Manager/function"
)

// CreateFunction creates a function on lambda
func (m *Manager) CreateFunction(funcName string, dir string, e env.Env, memoryS string, timeoutS string) (error) {
	var err error
	memorySize, err := strconv.ParseInt(memoryS, 10, 64)
	if err != nil {
		return err
	}
	timeout, err := strconv.ParseInt(timeoutS, 10, 64)
	if err != nil {
		return err
	}
	runtime := ""
	if e == env.PYTHON3 {
		runtime = "python3"
	} else if e == env.JAVA8 {
		runtime = "java8"
	} else {
		return errors.New("Not support env")
	}

	codeZip, err := os.Open(path.Join(dir, "code.zip"))
	if err != nil {
		return err
	}
	code, err := ioutil.ReadAll(codeZip)
	if err != nil {
		return err
	}
	m.storage.SetFunc(&SetFuncInput{
		FuncName: funcName,
		MemorySize: memorySize,
		Timeout: timeout,
		Runtime: runtime,
		Image: "",
	})
	m.storage.SetCode(funcName, code)

	return nil
}

// GetFunction fetches the metadata of a function
func (m *Manager) GetFunction(funcName string) (*function.Meta, error) {
	output, err := m.storage.GetMeta(funcName)
	if err != nil {
		return nil, err
	}
	return output, nil	
}

// InvokeFunction invokes a function by name
func (m *Manager) InvokeFunction(funcName string, args []byte) ([]byte, error) {
	return nil, errors.New("No native Serverless")
}

// ListFunction displays the functions
func (m *Manager) ListFunction() ([]*function.Meta, error) {
	functions := make([]*function.Meta, 0)
	// TODO
	return functions, nil
}

// DeleteFunction deletes a function
func (m *Manager) DeleteFunction(funcName string) (error) {
	// TODO
	return nil
}