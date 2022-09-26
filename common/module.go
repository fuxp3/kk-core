package common

import (
	"kk-core/core/module"
	"kk-core/core/task"
	"time"
)

type RunModule struct {
	module.BaseTaskModule
	Parallel bool
	Retry    int
	Hostname string
	Command  string
}

func (r *RunModule) Init() {
	r.Name = "执行命令Module"
	r.Desc = "在目标主机上执行命令"

	var timeout int64
	for _, v := range r.Runtime.GetAllHosts() {
		timeout += v.GetTimeout()
	}

	test := &task.RemoteTask{
		Name:     "执行命令Task",
		Desc:     "在目标主机上执行命令",
		Hosts:    r.Runtime.GetHostsByName(r.Hostname),
		Action:   &RunTask{Command: r.Command},
		Parallel: r.Parallel,
		Timeout:  time.Duration(timeout) * time.Second,
		Retry:    r.Retry,
	}

	r.Tasks = []task.Interface{
		test,
	}

}
