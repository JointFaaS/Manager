package httpmanager

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/JointFaaS/Manager/env"
)

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func deCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

// UploadHandler creates a new function
func (m* Manager) UploadHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	var funcName string
	var e env.Env
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			if part.FormName() == "funcName" {
				data, _ := ioutil.ReadAll(part)
				funcName = string(data)
			} else if part.FormName() == "env" {
				data, _ := ioutil.ReadAll(part)
				e, err = env.ConvertStrToEnv(string(data))
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				http.Error(w, "Unkown form data " + part.FormName() , http.StatusBadRequest)
				return
			}
		} else if part.FileName() == "code.zip" {
			dst, _ := os.Create(path.Join(dir, part.FileName()))
			defer dst.Close()
			io.Copy(dst, part)
		} else {
			http.Error(w, "Required code.zip missed", http.StatusBadRequest)
			return
		}
	}

	deCompress(path.Join(dir, "code.zip"), dir)

	err = m.platformManager.CreateFunction(funcName, dir, e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

