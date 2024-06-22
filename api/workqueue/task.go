package workqueue

import (
	"fmt"
	"os/exec"
)

type Task struct {
	Command    string
	Args       []string
	Type       interface{} // downlowd resource
	Name       string
	PreRunFunc func(string) error
	ResultFunc func(RunningStatus) (string, error)
	Dir        string
	Error      error
	Pid        string
}

type RunningStatus struct {
	IfSuccess  bool
	NodeStatus []ExecuteExitNodeStatus
	Pid        string
	ExecuteDir string
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

func (*Task) Execute(ts *Task) {
	fmt.Printf("执行任务", ts.Name)
	exec.Command(ts.Command, ts.Args...).Run()

}
