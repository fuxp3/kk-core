package test

import (
	"kk-core/core/module"
	"kk-core/core/task"
	"time"
)

type Module struct {
	module.BaseTaskModule
	Hostname string
	Command  string
}

func (h *Module) Init() {
	h.Name = "TestModule"
	h.Desc = "Test"

	var timeout int64
	for _, v := range h.Runtime.GetAllHosts() {
		timeout += v.GetTimeout()
	}

	test := &task.RemoteTask{
		Name:     "run",
		Desc:     "run",
		Hosts:    h.Runtime.GetHostsByName(h.Hostname),
		Action:   &Task{Command: h.Command},
		Parallel: true,
		Timeout:  time.Duration(timeout) * time.Second,
	}

	h.Tasks = []task.Interface{
		test,
	}

}
