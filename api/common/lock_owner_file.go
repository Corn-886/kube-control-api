package common

import (
	"errors"
	"github.com/sirupsen/logrus"
	"kube-control-api/config/constants"
	"os"
	"path/filepath"
	"sync"
)

/**
通过对 操作历史文件加锁，防止其他线程同时操作，并且记录
sample: lock  -->   /data/download/Resource/inventory.lastrun
*/
func LockOwner(owner_type string, owner_name string, rwMurex *sync.RWMutex) (*os.File, error) {
	lockFilePath := filepath.Join(constants.GET_DATA_DIR(), owner_type, owner_name, "inventory.lastrun")

	logrus.Trace("lockFilePath: ", lockFilePath)
	_, error := OpenOrCreateFile(lockFilePath)
	if error != nil {
		return nil, errors.New("Cannot create file " + lockFilePath + " : " + error.Error())
	}

	lockedFile, err := os.OpenFile(lockFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // 如果不存在则创建
	if err != nil {
		return nil, errors.New("Cannot open file " + lockFilePath + " : " + err.Error())
	}
	rwMurex.Lock()
	return lockedFile, nil
}

func UnLockOwner(rwMurex *sync.RWMutex) {
	rwMurex.Unlock()
}
