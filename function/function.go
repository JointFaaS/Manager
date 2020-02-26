package function

// Meta defines the metadata of a function
type Meta struct {
	FunctionName          string           `json:"functionName"`
	Description           string           `json:"description"`
	Runtime               string           `json:"runtime"`
	Handler               string           `json:"handler"`
	Timeout               int32            `json:"timeout"`
	MemorySize            int32            `json:"memorySize"`
	CodeSize              int64            `json:"codeSize"`
	CodeChecksum          string           `json:"codeChecksum"`
	EnvironmentVariables  map[string]string `json:"environmentVariables"`
	CreatedTime           string           `json:"createdTime"`
}