package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/JointFaaS/Manager/httpmanager"
	"gopkg.in/yaml.v2"
)

func logInit() {
	log.SetPrefix("TRACE: ")
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	logInit()
	var config httpmanager.Config
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cfgFile, err := ioutil.ReadFile(path.Join(home, "/.jfManager/config.yml"))
	if err != nil {
		panic(err)
	}

	err = yaml.UnmarshalStrict(cfgFile, &config)
	if err != nil {
		panic(err)
	}

	httpManager, err := httpmanager.NewManager(config)
	if err != nil {
		panic(err)
	}

	httpManager.SetRouter()

	log.Print("start listening")
    log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}