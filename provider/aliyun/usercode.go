package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// GetCodeURI returns codeURI of a given function
func (m *Manager) GetCodeURI(funcName string) (string, error) {
	codeOut, err := m.userCodeBucket.SignURL(funcName, oss.HTTPGet, 99999999)
	if err != nil {
		return "", err
	}
	return codeOut, nil
}

// SaveCode saves user's code
func (m *Manager) SaveCode(funcName string, file string) (error) {
	err := m.userCodeBucket.PutObjectFromFile(funcName, file)
	if err != nil {
		return err
	}
	return nil
}