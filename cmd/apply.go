package cmd

import (
	"github.com/spf13/cobra"
	"kk-core/core/connector"
	"kk-core/core/module"
	"kk-core/core/pipeline"
	"kk-core/test"
)

var (
	name    string
	command string
	runCmd  = &cobra.Command{
		Use:   "run",
		Short: "kk-core",
		Long:  "测试kk-core任务调度框架",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		Run: func(cmd *cobra.Command, args []string) {
			p := pipeline.Pipeline{
				Name: "测试流水线",
				Modules: []module.Module{
					&test.Module{
						Hostname: name,
						Command:  command,
					},
				},
				Runtime: NewRuntime(),
			}

			p.Start()
		},
	}
)

func init() {
	runCmd.Flags().StringVarP(&name, "name", "n", "", "主机名")
	runCmd.Flags().StringVarP(&command, "command", "c", "", "命令")
	runCmd.MarkFlagsRequiredTogether("name", "command")
}

func NewRuntime() connector.Runtime {
	base := connector.NewBaseRuntime("test", connector.NewDialer(), false, false)
	hosts := []connector.Host{
		&connector.BaseHost{
			Name:            "master2",
			User:            "root",
			Port:            22,
			Password:        "123",
			InternalAddress: "192.168.108.132",
			Address:         "192.168.108.132",
		},
		&connector.BaseHost{
			Name:            "control",
			User:            "root",
			Port:            22,
			Password:        "123",
			InternalAddress: "192.168.108.143",
			Address:         "192.168.108.143",
		},
	}
	base.SetAllHosts(hosts)
	return &base
}
