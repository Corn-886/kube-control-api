package oscmd

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"kube-control-api/api/common"
	"kube-control-api/config/constants"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var runningProcesses = map[string]*os.Process{}

func (runner *CommandRunnerConfig) ToString(execute_dir_path string, pid string) string {
	result := "---"
	result += "\nowner_type: " + runner.Type
	result += "\nowner_name: " + runner.Name
	result += "\nhistory: " + execute_dir_path
	result += "\npid: " + pid
	result += "\ndir: " + runner.Dir
	result += "\ncmd: " + runner.Command
	result += "\nargs:"
	for _, arg := range runner.Args {
		result += "\n  - " + arg
	}
	result += "\nenv:"

	result += "\nshell: cd " + runner.Dir

	result += " && " + runner.Command
	for _, arg := range runner.Args {
		result += " " + arg
	}
	return result
}

type SuccessTask struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"succeedTime"`
	Pid       string    `json:"pid"`
}

type CommandRunnerConfig struct {
	Command    string
	Args       []string
	mutex      *sync.Mutex //互斥锁
	Type       string      // 用于定位上锁文件，不能修改
	Name       string      // 用于定位上锁文件，不能修改
	PreRunFunc func(string) error
	ResultFunc func(RunningStatus) (string, error)
	Dir        string
	Error      error
	Pid        string
}
type ExecuteExitNodeStatus struct {
	NodeName    string
	OK          string
	Changed     string
	Unreachable string
	Failed      string
	Skipped     string
	Rescued     string
	Ignored     string
}

type RunningStatus struct {
	IfSuccess  bool
	NodeStatus []ExecuteExitNodeStatus
	Pid        string
	ExecuteDir string
}

func parseAnsibleRecapLine(line string) ExecuteExitNodeStatus {
	result := make(map[string]string)
	s1 := strings.Split(line, ":")

	result["node"] = strings.Trim(s1[0], " ")
	s2 := strings.Split(s1[1], " ")
	for _, kv := range s2 {
		s3 := strings.Split(kv, "=")
		if len(s3) == 1 {
			continue
		}
		k := strings.Trim(s3[0], " ")
		v := strings.Trim(s3[1], " ")
		result[k] = v
	}

	nodeStatus := ExecuteExitNodeStatus{
		NodeName:    trimColoredText(result["node"]),
		OK:          result["ok"],
		Changed:     trimColoredText(result["changed"]),
		Unreachable: result["unreachable"],
		Failed:      result["failed"],
		Skipped:     result["skipped"],
		Rescued:     result["rescured"],
		Ignored:     result["ignored"],
	}
	return nodeStatus
}

func trimColoredText(text string) string {
	result := text
	if strings.Index(result, "[0;32m") > 0 {
		result = result[strings.Index(result, "[0;32m")+6:]
	}
	if strings.Index(result, "[0;33m") > 0 {
		result = result[strings.Index(result, "[0;33m")+6:]
	}
	if strings.Index(result, "\033[0m") > 0 {
		result = result[:strings.Index(result, "\033[0m")]
	}
	return result
}

func (runner *CommandRunnerConfig) RunCommand() error {
	runner.mutex = &sync.Mutex{}
	runner.mutex.Lock() // 第一次加锁，确保go runcommand 能运行完
	go runner.runCommand()
	logrus.Trace("geting lock in main thread...")
	runner.mutex.Lock()
	logrus.Trace("Got response from exec: ", runner.Error)
	runner.mutex.Unlock()
	return runner.Error
}

func (runner *CommandRunnerConfig) runCommand() {
	defer func() {
		if err := recover(); err != nil {
			println(err)
			runner.Error = errors.New("unknow error ")
			runner.mutex.Unlock()
		}
		return
	}()

	lockFile, lockedFile, err := common.LockOwner(runner.Type, runner.Name)

	if err != nil {
		runner.Error = err
		runner.mutex.Unlock()
		return
	}

	common.UnLockOwner(lockFile)

	pid := time.Now().Format("2006-01-02_15-04-05.999") + "_" + runner.Type
	historyPath := filepath.Join(constants.GET_DATA_DIR(), runner.Type, runner.Name, "history")
	if err := common.CreateDirIfNotExists(historyPath); err != nil {
		runner.Error = errors.New("cannot create historyDir : " + historyPath + " : " + err.Error())
		runner.mutex.Unlock()
		return
	}
	execute_dir_path := filepath.Join(historyPath, pid)
	if err := common.CreateDirIfNotExists(execute_dir_path); err != nil {
		runner.Error = errors.New("cannot create runDir : " + execute_dir_path + " : " + err.Error())
		runner.mutex.Unlock()
		return
	}

	logFilePath := filepath.Join(execute_dir_path, "execute.log")
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		runner.Error = errors.New("cannot create logFile : " + logFilePath + " : " + err.Error())
		runner.mutex.Unlock()
		return
	}
	defer logFile.Sync()
	defer logFile.Close() //确保所有文件都释放

	cmd := exec.Command(runner.Command, runner.Args...)

	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Dir = runner.Dir

	if err := cmd.Start(); err != nil {
		runner.Error = errors.New("failed to start command " + cmd.String() + " : " + err.Error())
		os.Remove(execute_dir_path)
		runner.mutex.Unlock()
		return
	}

	runningProcesses[pid] = cmd.Process

	logrus.Trace("started command " + cmd.String())
	ioutil.WriteFile(filepath.Join(execute_dir_path, "execute.command"), []byte(cmd.String()), 0666)
	ioutil.WriteFile(filepath.Join(execute_dir_path, "execute.yaml"), []byte(runner.ToString(execute_dir_path, pid)), 0666)

	if err := lockedFile.Truncate(0); err != nil {
		runner.Error = errors.New("failed to truncate lockFile : " + err.Error())
	}

	_, err = lockedFile.WriteString(pid)
	if err != nil {
		runner.Error = errors.New("failed to write pid " + err.Error())
	}

	runner.Pid = pid

	runner.mutex.Unlock()
	cmd.Wait()

	delete(runningProcesses, pid)

	if runner.ResultFunc != nil {
		logFile.WriteString("\n\n\nKUBOARD SPRAY *****************************************************************\n")
		logs, err := ioutil.ReadFile(logFilePath)
		if err != nil {
			return
		}
		logStr := string(logs)

		if strings.LastIndex(logStr, "PLAY RECAP *********************************************************************") < 0 {
			logFile.WriteString("\033[31m\033[01m\033[05m  No PLAY RECAP found.\033[0m\n")
			logFile.WriteString("\033[31m\033[01m\033[05m  执行出错。\033[0m\n")
			logrus.Warn("No ansbile-playbook recap.")
			return
		}
		recap := logStr[strings.LastIndex(logStr, "PLAY RECAP *********************************************************************")+81:]
		recap = recap[:strings.Index(recap, "\n\n")]

		lines := strings.Split(recap, "\n")
		status := []ExecuteExitNodeStatus{}

		for _, line := range lines { // 捕捉关键字，如果日志存在failed/unreachable 则返回false
			if strings.Index(line, "failed=") > 0 && strings.Index(line, "unreachable=") > 0 && strings.Index(line, "ok=") > 0 {
				status = append(status, parseAnsibleRecapLine(line))
			}
		}

		ifSuccess := len(status) > 0
		for _, nodestatus := range status {
			if nodestatus.Unreachable != "0" || nodestatus.Failed != "0" {
				ifSuccess = false
			}
		}

		exitStatus := RunningStatus{
			IfSuccess:  ifSuccess,
			NodeStatus: status,
			Pid:        pid,
			ExecuteDir: execute_dir_path,
		}
		message, err := runner.ResultFunc(exitStatus)
		if err != nil {
			logFile.WriteString("Error in running result: " + err.Error() + "\n")
		}

		if ifSuccess {
			task := common.SuccessTask{
				Type:      runner.Type,
				Timestamp: time.Now(),
				Pid:       pid,
			}
			if err := common.AddSuccessTask(runner.Type, runner.Name, task); err != nil {
				logrus.Warn("failed to add success task: ", err)
			}
		}

		logFile.WriteString(message)
	}
}
