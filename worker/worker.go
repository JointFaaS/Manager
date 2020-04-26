package worker

import (
	"context"
	"errors"
	"log"

	wpb "github.com/JointFaaS/Manager/pb/worker"
	"google.golang.org/grpc"
)

// Worker is the wrapper handler fo interacting with worker
type Worker struct {
	addr string
	id string

	wc wpb.WorkerClient
	activeFunctions map[string]bool // the number of the specified function instances
}

// New creates a worker handler
func New(addr string, id string) (*Worker, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("can not connect with server %v", err)
		return nil, err
	}
	rpcClient := wpb.NewWorkerClient(conn)
	return &Worker{
		addr: addr,
		id: id,
		wc: rpcClient,
		activeFunctions: make(map[string]bool),
	}, nil
}

// InitFunction initialise an instance of the given function
func (w *Worker) InitFunction(ctx context.Context, funcName string, image string, codeURI string) (error) {
	res, err := w.wc.InitFunction(ctx, &wpb.InitFunctionRequest{
		FuncName: funcName,
		Image: image,
		CodeURI: codeURI,
		Runtime: "",
		Timeout: 3,
		MemorySize: 128,
	})
	if err != nil {
		return err
	}
	if res.GetCode() != wpb.InitFunctionResponse_OK {
		return errors.New(res.GetMsg())
	}
	w.activeFunctions[funcName] = true

	return nil
}

// CallFunction 
func (w *Worker) CallFunction(ctx context.Context, funcName string, args []byte) ([]byte, error){
	res, err := w.wc.Invoke(ctx, &wpb.InvokeRequest{
		Name: funcName,
		Payload: args,
	})
	if err != nil {
		return nil, err
	}
	if res.GetCode() != wpb.InvokeResponse_OK {
		return nil, errors.New(res.GetCode().String())
	}
	return res.GetOutput(), nil
}

// HasFunction
func (w *Worker) HasFunction(funcName string) (bool){
	e, isPresent := w.activeFunctions[funcName]
	return isPresent && e
}