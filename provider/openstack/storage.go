package openstack

import (
	"github.com/JointFaaS/Manager/function"
)

type SetFuncInput struct {
	FuncName string
	Timeout int64
	MemorySize int64
	Runtime string
	Image string
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

func (s *LocalFileStorage) SetFunc(*SetFuncInput) error {
	return nil
}

func (s *LocalFileStorage) SetCode(string, []byte) (error) {
	return nil
}

func (s *LocalFileStorage) GetMeta(string) (*function.Meta, error) {
	return nil, nil
}

func (s *LocalFileStorage) GetCode(string) ([]byte, error) {
	return nil, nil
}

func (s *LocalFileStorage) Del(string) (error) {
	return nil
}