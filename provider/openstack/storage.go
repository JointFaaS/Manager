package openstack

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/JointFaaS/Manager/function"
)

func ensureFuncDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModeDir)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if info.IsDir() == false {
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
		err = os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return err
		}
	}
	return nil
}

type SetFuncInput struct {
	FuncName   string
	Timeout    int64
	MemorySize int64
	Runtime    string
	Image      string
}

type Storage interface {
	SetFunc(*SetFuncInput) error
	SetCode(string, []byte) error
	GetMeta(string) (*function.Meta, error)
	GetCode(string) ([]byte, error)
	Del(string) error
}

type LocalFileStorage struct {
	rootDir string
}

func NewLocalFileStorage(rootDir string) *LocalFileStorage {
	return &LocalFileStorage{
		rootDir: rootDir,
	}
}

func (s *LocalFileStorage) SetFunc(input *SetFuncInput) error {
	dir := path.Join(s.rootDir, input.FuncName)
	err := ensureFuncDir(dir)
	if err != nil {
		return err
	}
	metaPath := path.Join(dir, "meta")
	metaFile, err := os.OpenFile(metaPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	m := new(function.Meta)
	m.MemorySize = input.MemorySize
	m.FunctionName = input.FuncName
	m.Runtime = input.Runtime
	m.Timeout = input.Timeout
	json.NewEncoder(metaFile).Encode(m)
	return nil
}

func (s *LocalFileStorage) SetCode(funcName string, code []byte) error {
	dir := path.Join(s.rootDir, funcName)
	err := ensureFuncDir(dir)
	if err != nil {
		return err
	}
	codePath := path.Join(dir, "code")
	codeFile, err := os.OpenFile(codePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer codeFile.Close()
	codeFile.Write(code)
	return nil
}

func (s *LocalFileStorage) GetMeta(funcName string) (*function.Meta, error) {
	metaPath := path.Join(s.rootDir, funcName, "meta")
	metaFile, err := os.OpenFile(metaPath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer metaFile.Close()
	m := new(function.Meta)
	err = json.NewDecoder(metaFile).Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *LocalFileStorage) GetCode(funcName string) ([]byte, error) {
	codePath := path.Join(s.rootDir, funcName, "code")
	codeFile, err := os.OpenFile(codePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer codeFile.Close()
	codeBytes, err := ioutil.ReadAll(codeFile)
	if err != nil {
		return nil, err
	}
	return codeBytes, nil
}

func (s *LocalFileStorage) Del(funcName string) error {
	os.RemoveAll(path.Join(s.rootDir, funcName))
	return nil
}
