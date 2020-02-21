package aliyun

import (
	"os"
	"path"
)

// inject neccessary index.handler adaptor for aliyun function
func (m* Manager) createPython3Function(dir string) error {
	index, _ := os.Create(path.Join(dir, "index.py"))

}