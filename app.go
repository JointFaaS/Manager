package main

import (
	"log"
	"net/http"

	"github.com/JointFaaS/Manager/controller"
)

func logInit() {
	log.SetPrefix("TRACE: ")
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	logInit()
	var config controller.Config
	_, err := controller.NewManager(config)
	if err != nil {
		panic(err)
	}
	
	createFunctionHandler := func (w http.ResponseWriter, r *http.Request) {

	}
	invokeHandler := func (w http.ResponseWriter, r *http.Request) {

	}

	http.HandleFunc("/createfunction", createFunctionHandler)
	http.HandleFunc("/invokefunction/*", invokeHandler)

	log.Print("start listening")
    log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}