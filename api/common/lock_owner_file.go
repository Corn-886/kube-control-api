package common

import (
	"errors"
	"github.com/nightlyone/lockfile"
	"github.com/sirupsen/logrus"
	"kube-control-api/config/constants"
	"os"
	"path/filepath"
)

/**
通过对 操作历史文件加锁，防止其他线程同时操作，并且记录
sample: lock  -->   /data/download/Resource/inventory.lastrun
*/
func LockOwner(owner_type string, owner_name string) (*lockfile.Lockfile, *os.File, error) {
	lockFilePath := filepath.Join(constants.GET_DATA_DIR(), owner_type, owner_name, "inventory.lastrun")

	logrus.Trace("lockFilePath: ", lockFilePath)
	_, error := OpenOrCreateFile(lockFilePath)
	if error != nil {
		return nil, nil, errors.New("Cannot create file " + lockFilePath + " : " + error.Error())
	}

	lockedFile, err := os.OpenFile(lockFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // 如果不存在则创建
	if err != nil {
		return nil, nil, errors.New("Cannot open file " + lockFilePath + " : " + err.Error())
	}

	lockOwnerFile, err := lockfile.New(lockFilePath)

	if err != nil {
		return nil, nil, errors.New("Cannot init lock " + lockFilePath + " : " + err.Error())
	}

	err = lockOwnerFile.TryLock()
	if err != nil {
		return nil, nil, errors.New("Cannot  lock file" + lockFilePath + " : " + err.Error())
	}

	return &lockOwnerFile, lockedFile, nil
}

func UnLockOwner(lockedfile *lockfile.Lockfile) {
	lockedfile.Unlock()
}
