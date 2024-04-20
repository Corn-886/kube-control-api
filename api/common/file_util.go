package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"kube-control-api/config/constants"
	"os"
	"path/filepath"
	"time"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func CreateDirIfNotExists(path string) error {
	if PathExists(path) {
		return nil
	}
	return os.Mkdir(path, os.ModePerm)
}
func OpenOrCreateFile(lockFilePath string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(lockFilePath), os.ModePerm)
	if err != nil {
		return nil, err
	}

	lockedFile, err := os.OpenFile(lockFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return lockedFile, nil
}

type SuccessTask struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"succeedTime"`
	Pid       string    `json:"pid"`
}

type SuccessTasks []SuccessTask

func AddSuccessTask(ownerType string, ownerName string, task SuccessTask) error {
	tasks, err := ReadSuccessTasks(ownerType, ownerName)
	if err != nil {
		return err
	}
	tasks = append(tasks, task)
	content, err := json.Marshal(tasks)
	if err != nil {
		return errors.New("failed tp marshal tasks: " + err.Error())
	}
	return ioutil.WriteFile(successFilePath(ownerType, ownerName), content, 0666)
}

func ReadSuccessTasks(ownerType string, ownerName string) (SuccessTasks, error) {
	var tasks SuccessTasks
	content, err := ioutil.ReadFile(successFilePath(ownerType, ownerName))
	if err != nil {
		content = []byte("[]")
	}
	if err := json.Unmarshal(content, &tasks); err != nil {
		return nil, errors.New("failed to unmarshal file " + successFilePath(ownerType, ownerName) + " : " + err.Error())
	}

	return tasks, nil
}

func successFilePath(ownerType string, ownerName string) string {
	return filepath.Join(constants.GET_DATA_DIR(), ownerType, ownerName, "success.json")
}
