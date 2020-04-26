package openstack

import (
	"io/ioutil"
	"os"
)

func (m *Manager) GetCodeURI(funcName string) (string, error) {
	return "", nil
}

func (m *Manager) GetImage(funcName string) (string, error) {
	return "", nil
}

func (m *Manager) SaveCode(funcName string, file string) (error) {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = m.storage.SetCode(funcName, body)
	if err != nil {
		return err
	}
	return nil
}