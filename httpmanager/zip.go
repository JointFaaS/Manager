package httpmanager

import (
    "archive/zip"
    "io"
    "path"
    "os"
)

func deCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

	defer reader.Close()
	for _, innerFile := range reader.File {
        info := innerFile.FileInfo()
        if info.IsDir() {
            err = os.MkdirAll(path.Join(dest, innerFile.Name), os.ModePerm)
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