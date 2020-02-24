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
	for _, innerFile := range reader.File {
        info := innerFile.FileInfo()
        if info.IsDir() {
            err = os.MkdirAll(innerFile.Name, os.ModePerm)
            if err != nil {
                return err
            }
            continue
        }
        srcFile, err := innerFile.Open()
        if err != nil {
            return err
        }
        defer srcFile.Close()
        newFile, err := os.Create(path.Join(dest, innerFile.Name))
        if err != nil {
			return err
        }
        io.Copy(newFile, srcFile)
        newFile.Close()
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

	err = m.platformManager.CreateFunction(funcName, dir, e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

