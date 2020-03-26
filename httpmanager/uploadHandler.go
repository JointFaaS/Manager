package httpmanager

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/JointFaaS/Manager/env"
)

type uploadBody struct {
	FuncName string `json:"funcName"`
	CodeZip  string `json:"codeZip"`
	Env string `json:"env"`
	MemorySize     string `json:"memorySize"`
	Timeout string `json:"timeout"`
}

// UploadHandler creates a new function
func (m* Manager) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not support method", http.StatusBadRequest)
		return
	}
	var req uploadBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	funcName := req.FuncName
	codeZip, err := base64.StdEncoding.DecodeString(req.CodeZip)
	if err != nil {
		http.Error(w, "Fail to read CodeZip", http.StatusBadRequest)
		return
	}
	e, err := env.ConvertStrToEnv(req.Env)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	dst, _ := os.Create(path.Join(dir, "code.zip"))
	defer dst.Close()
	io.Copy(dst, bytes.NewBuffer(codeZip))

	err = m.platformManager.SaveCode(funcName, path.Join(dir, "code.zip"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = os.Mkdir(path.Join(dir, "code"), os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = deCompress(path.Join(dir, "code.zip"), path.Join(dir, "code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.platformManager.CreateFunction(funcName, dir, e, req.MemorySize, req.Timeout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

