package prehello

import (
	"kk-core/core/action"
	"kk-core/core/connector"
	"kk-core/core/logger"
)

type HelloTask struct {
	action.BaseAction
}

func (h *HelloTask) Execute(runtime connector.Runtime) error {
	hello, err := runtime.GetRunner().SudoCmd("echo 'Greetings, KubeKey!'", false)
	if err != nil {
		return err
	}
	logger.Log.Messagef(runtime.RemoteHost().GetName(), hello)
	return nil
}
