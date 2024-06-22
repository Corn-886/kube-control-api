package oscmd

import (
	"os"
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
