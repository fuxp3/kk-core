package common

import (
	"kk-core/core/action"
	"kk-core/core/connector"
	"kk-core/core/logger"
)

type RunTask struct {
	action.BaseAction
	Command string
}

func (r *RunTask) Execute(runtime connector.Runtime) error {
	//hello, err := runtime.GetRunner().SudoCmd("echo 'Greetings, KubeKey!'", true)
	hello, err := runtime.GetRunner().SudoCmd(r.Command, false)
	if err != nil {
		return err
	}
	logger.Log.Messagef(runtime.RemoteHost().GetName(), hello)
	return nil
}
