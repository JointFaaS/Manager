package worker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Worker is the wrapper handler fo interacting with worker
type Worker struct {
	addr string
	id string

	activeFunctions map[string]int32 // the number of the specified function instances
}

type initRequestBody struct {
	FuncName string `json:"funcName"`
	Image string	`json:"image"`
	CodeURI string	`json:"codeURI"`
}

type callRequestBody struct {
	FuncName string `json:"funcName"`
	Args string `json:"args"`
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
func (w *Worker) InitFunction(funcName string, image string, codeURI string) (error) {
	body := initRequestBody{
		FuncName: funcName,
		Image: image,
		CodeURI: codeURI,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = http.Post("http://" + w.addr + "/init", "application/json;charset=UTF-8", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	return nil
}

// CallFunction 
func (w *Worker) CallFunction(funcName string, args []byte) ([]byte, error){
	resp, err := http.Post("http://" + w.addr + "/call?funcName=" + funcName, "application/json;charset=UTF-8", bytes.NewReader(args))
	if err != nil {
		return nil, err
	}
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ret, nil
}