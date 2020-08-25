package aws

import (
	"path"
	"os"
)

func (m *Manager) injectPython3Handler(dir string) (error) {
	err := os.Link(path.Join(m.addonsDir, "python3", "jointfaas.py") ,path.Join(dir, "jointfaas.py"))
	if err != nil {
		return err
	}
	return nil
}