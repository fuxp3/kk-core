package cmd

import (
	"github.com/spf13/cobra"
	"kk-core/core/connector"
	"kk-core/core/module"
	"kk-core/core/pipeline"
	"kk-core/prehello"
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "测试kk-core任务调度框架",
	Long:  "测试kk-core任务调度框架",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		m := []module.Module{
			&prehello.HelloModule{},
		}

		runtime := NewRuntime()

		p := pipeline.Pipeline{
			Name:    "我的测试流水线",
			Modules: m,
			Runtime: runtime,
		}

		p.Start()
	},
}

func NewRuntime() connector.Runtime {
	base := connector.NewBaseRuntime("test", connector.NewDialer(), false, false)
	return &base
}
