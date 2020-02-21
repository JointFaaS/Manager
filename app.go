package main

import (
	"io/ioutil"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/JointFaaS/Manager/controller"
	"github.com/JointFaaS/Manager/env"	
)

func logInit() {
	log.SetPrefix("TRACE: ")
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	logInit()
	var config controller.Config
	manager, err := controller.NewManager(config)
	if err != nil {
		panic(err)
	}

	createFunctionHandler := func (w http.ResponseWriter, r *http.Request) {
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		formData := make(map[string][]byte)
		dir, err := ioutil.TempDir("", "")
		defer os.RemoveAll(dir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if part.FileName() == "" {
				data, _ := ioutil.ReadAll(part)
				formData[part.FormName()] = data
			} else{
				dst, _ := os.Create(path.Join(dir, part.FileName()))
				defer dst.Close()
				io.Copy(dst, part)
			}
		}
		err = manager.UploadFunction(string(formData["funcName"]), dir, env.Env(formData["env"]))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	invokeHandler := func (w http.ResponseWriter, r *http.Request) {

	}

	http.HandleFunc("/createfunction", createFunctionHandler)
	http.HandleFunc("/invokefunction/", invokeHandler)

	log.Print("start listening")
    log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}