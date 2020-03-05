package aliyun

import (
	"os"
	"path"
)

// inject neccessary AliIndex.class for aliyun function
func (m* Manager) createJava8Function(dir string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// TODO: config
	err = os.Link(path.Join(home, ".jfManager", "ali", "java8", "AliIndex.class") ,path.Join(dir, "jointfaas", "AliIndex.class"))
	if err != nil {
		return err
	}
	return nil
}