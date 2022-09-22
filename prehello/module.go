package prehello

import (
	"kk-core/core/module"
	"kk-core/core/task"
	"time"
)

type HelloModule struct {
	module.BaseTaskModule
}

func (h *HelloModule) Init() {
	h.Name = "GreetingsModule"
	h.Desc = "Greetings"

	var timeout int64
	for _, v := range h.Runtime.GetAllHosts() {
		timeout += v.GetTimeout()
	}

	hello := &task.RemoteTask{
		Name:     "Greetings",
		Desc:     "Greetings",
		Hosts:    h.Runtime.GetAllHosts(),
		Action:   new(HelloTask),
		Parallel: true,
		Timeout:  time.Duration(timeout) * time.Second,
	}

	h.Tasks = []task.Interface{
		hello,
	}

}
