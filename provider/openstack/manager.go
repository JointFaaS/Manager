package openstack

import (
)

// Manager exposes the interfaces about resource CRUD on aws
type Manager struct {
	storage Storage
}

// Config defines the neccessary settings about openstack.Manager
type Config struct {
	StorageRootDir string `yaml:"storageRootDir"`
	user string `yaml:"user"`
	secret string `yaml:"secret"`
	project string `yaml:"project"`
	zone string `yaml:"zone"`
	// TODO: support auto-scale
}

// NewManagerWithConfig create a manager with given configuration
func NewManagerWithConfig(config Config) (*Manager, error){
	manager := &Manager{
		storage: NewLocalFileStorage(config.StorageRootDir),
	}
	return manager, nil
}