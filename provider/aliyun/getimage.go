package aliyun

// GetImage returns image of a function
func (m *Manager) GetImage(funcName string) (string, error) {
	// TODO: support more envs
	return "registry.cn-shanghai.aliyuncs.com/veia/hcloud-py:latest", nil
}