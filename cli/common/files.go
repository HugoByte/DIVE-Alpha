package common

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type diveFileHandler struct{}

func NewDiveFileHandler() *diveFileHandler {
	return &diveFileHandler{}
}
func (df *diveFileHandler) ReadFile(filePath string) ([]byte, error) {

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, Errorcf(FileError, "Error While Reading File %s", err.Error())
	}

	return fileData, nil
}

func (df *diveFileHandler) ReadJson(fileName string, obj interface{}) error {
	pwd, err := df.GetPwd()
	if err != nil {
		return WrapMessageToError(err, "Error While Reading File")
	}

	filePath := filepath.Join(pwd, fileName)

	data, err := df.ReadFile(filePath)
	if err != nil {
		return err
	}

	if len(data) != 0 {
		if err := json.Unmarshal(data, obj); err != nil {
			return WrapCodeToError(err, FileError, "Failed to Unmarshal Data")
		}
	}

	return nil
}
func (df *diveFileHandler) ReadAppFile(fileName string) ([]byte, error) {

	uhd, err := df.GetHomeDir()
	if err != nil {
		return nil, WrapMessageToError(err, "Failed to Read App File")
	}
	appFilePath := filepath.Join(uhd, ".dive", fileName)

	data, err := df.ReadFile(appFilePath)

	if err != nil {
		return nil, WrapMessageToError(err, "Failed To Read App File")
	}

	return data, nil
}

func (df *diveFileHandler) WriteAppFile(fileName string, data []byte) error {

	uhd, err := df.GetHomeDir()
	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Write App File %s", fileName)
	}
	appFileDir := filepath.Join(uhd, ".dive")

	err = df.MkdirAll(appFileDir, os.ModePerm)

	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Write App File %s", fileName)
	}

	appFilePath := filepath.Join(appFileDir, fileName)

	file, err := df.OpenFile(appFilePath, "append|write|create", 0644)
	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Write App File %s", fileName)
	}

	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Write App File %s", fileName)
	}

	return nil
}

func (df *diveFileHandler) WriteFile(fileName string, data []byte) error {

	pwd, err := df.GetPwd()

	if err != nil {
		return WrapMessageToErrorf(err, "Failed to Write File %s", fileName)
	}
	filePath := filepath.Join(pwd, fileName)

	file, err := df.OpenFile(filePath, "write|append|create", 0644)

	if err != nil {
		return WrapMessageToError(err, "Failed")
	}

	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Write App File %s", fileName)
	}

	return nil
}

func (df *diveFileHandler) WriteJson(fileName string, data interface{}) error {

	serializedData, err := json.Marshal(data)

	if err != nil {
		return WithCode(err, FileError)
	}

	err = df.WriteFile(fileName, serializedData)
	if err != nil {
		return WithCode(err, FileError)
	}
	return nil
}

func (df *diveFileHandler) GetPwd() (string, error) {

	pwd, err := os.Getwd()
	if err != nil {
		return "", Errorc(FileError, "Failed To Get PWD")
	}
	return pwd, err
}

func (df *diveFileHandler) MkdirAll(dirPath string, permission fs.FileMode) error {

	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, permission); err != nil {
			return WrapCodeToError(err, FileError, "Failed to Create Directory")
		}
	} else if err != nil {

		return WrapCodeToError(err, FileError, "Failed to check directory existence")
	}

	return nil
}

func (df *diveFileHandler) OpenFile(filePath string, fileOpenMode string, permission int) (*os.File, error) {
	mode := parseFileOpenMode(fileOpenMode)
	file, err := os.OpenFile(filePath, mode, fs.FileMode(permission))
	if err != nil {
		return nil, WrapCodeToError(err, FileError, "Failed to Open File")
	}

	return file, nil

}

func (df *diveFileHandler) GetHomeDir() (string, error) {

	uhd, err := os.UserHomeDir()
	if err != nil {
		return "", Errorc(FileError, "Failed To Get User HomeDir")
	}
	return uhd, err
}

func parseFileOpenMode(fileOpenMode string) int {
	modes := strings.Split(fileOpenMode, "|")

	var mode int
	for _, m := range modes {
		switch strings.TrimSpace(m) {
		case "append":
			mode |= os.O_APPEND
		case "create":
			mode |= os.O_CREATE
		case "truncate":
			mode |= os.O_TRUNC
		case "write":
			mode |= os.O_WRONLY
		case "readwrite":
			mode |= os.O_RDWR
		case "read":
			mode |= os.O_RDONLY
		}

	}

	return mode
}
