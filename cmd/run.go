package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"kk-core/common"
	"kk-core/core/connector"
	"kk-core/core/module"
	"kk-core/core/pipeline"
)

var (
	name    string
	command string
	hosts   []connector.BaseHost
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
					&common.Module{
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
	file, err := ioutil.ReadFile("./host.yaml")
	if err != nil {
		panic(any("未加载到主机配置文件"))
	}
	yaml.Unmarshal(file, &hosts)

	runCmd.Flags().StringVarP(&name, "name", "n", "", "主机名")
	runCmd.Flags().StringVarP(&command, "command", "c", "", "命令")
	runCmd.MarkFlagsRequiredTogether("name", "command")
}

func NewRuntime() connector.Runtime {
	base := connector.NewBaseRuntime("run", connector.NewDialer(), false, false)
	for _, h := range hosts {
		h := h
		base.AppendHost(&h)
	}
	return &base
}
