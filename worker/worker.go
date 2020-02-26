package worker

// Worker is the wrapper handler fo interacting with worker
type Worker struct {
	addr string
	id string

	activeFunctions map[string]int32 // the number of the specified function instances
}

// New creates a worker handler
func New(addr string, id string) (*Worker, error) {
	// TODO: validation
	return &Worker{
		addr: addr,
		id: id,
		activeFunctions: make(map[string]int32),
	}, nil
}

// InitFunction initialise an instance of the given function
func (*Worker) InitFunction(funcName string, image string, codeURI string) {

}

// CallFunction 
func (*Worker) CallFunction(funcName string, args []byte) ([]byte, error){
	return nil, nil
}