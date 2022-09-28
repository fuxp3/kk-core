/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"kk-core/common"
	"kk-core/core/connector"
	"kk-core/core/module"
	"kk-core/core/pipeline"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var (
	name     string
	command  string
	parallel bool
	retry    int
	hosts    []connector.BaseHost
	rootCmd  = &cobra.Command{
		Use:   "run",
		Short: "执行命令",
		Long:  "远程主机执行命令",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		Run: func(cmd *cobra.Command, args []string) {
			p := pipeline.Pipeline{
				Name: "执行命令流水线",
				Modules: []module.Module{
					&common.RunModule{
						Hostname: name,
						Command:  command,
						Parallel: parallel,
						Retry:    retry,
					},
				},
				Runtime: NewRuntime(),
			}

			p.Start()
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kk-core.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	file, err := ioutil.ReadFile("./host.yaml")
	if err != nil {
		panic(any("未加载到主机配置文件"))
	}
	yaml.Unmarshal(file, &hosts)

	rootCmd.Flags().StringVarP(&name, "name", "n", "", "主机,多个主机用\",\"分隔")
	rootCmd.Flags().StringVarP(&command, "command", "c", "", "命令")
	rootCmd.Flags().BoolVarP(&parallel, "parallel", "p", true, "是否并行")
	rootCmd.Flags().IntVarP(&retry, "retry", "r", 1, "重试次数")
	rootCmd.MarkFlagsRequiredTogether("name", "command")

	rootCmd.AddCommand(execCmd)
}

func NewRuntime() connector.Runtime {
	base := connector.NewBaseRuntime("run", connector.NewDialer(), false, false)
	for _, h := range hosts {
		h := h
		base.AppendHost(&h)
	}
	return &base
}
