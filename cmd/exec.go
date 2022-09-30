package cmd

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"kk-core/core/connector"
	"kk-core/core/logger"
	"os"
	"strings"
	"time"
)

var (
	execCmd = &cobra.Command{
		Use:   "exec",
		Short: "执行命令",
		Long:  "远程主机执行命令",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		Run: func(cmd *cobra.Command, args []string) {
			runtime := NewRuntime()

			if name == "" {
				name = strings.Join(args, ",")
			}

			host := GetHostsByName(name)
			if host == nil {
				logger.Log.Warnf("未找到主机：%s", name)
				return
			}

			/*err := executeCmd("echo hello", host)
			if err != nil {
				logger.Log.Errorln(err)
			}*/
			var command string

			_, err := executeCmd(runtime, "echo hello", host)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("连接主机...")
			fmt.Println("连接主机成功")
			fmt.Println("Last login:", time.Now())
			fmt.Printf("[root@%s ~]# ", host.GetName())

			reader := bufio.NewReader(os.Stdin) //终端输入
			for {
				line, _, _ := reader.ReadLine()
				command = string(line)
				/*for {
					b, err := reader.ReadByte()
					if err != nil {
						break
					}
					if b == byte('\r') {
						break
					}
					if b == byte('\n') {
						break
					}
					command += string(b)
				}*/

				if command != "" {
					respMsg, err := executeCmd(runtime, command, host)
					if err != nil {
						logger.Log.Errorln(err)
						//return
					}
					fmt.Println(respMsg)
					fmt.Printf("[%s@%s ~]# ", host.GetUser(), host.GetName())
				} else {
					fmt.Printf("[%s@%s ~]# ", host.GetUser(), host.GetName())
				}
				command = ""
			}
		},
	}
)

func init() {
	execCmd.Flags().StringVarP(&name, "name", "n", "", "主机,多个主机用\",\"分隔")
}

func GetHostsByName(name string) connector.Host {
	for _, h := range hosts {
		if h.GetName() == name {
			return &h
		}
	}
	return nil
}

func executeCmd(runtime connector.Runtime, cmd string, host connector.Host) (string, error) {
	conn, err := runtime.GetConnector().Connect(host)
	defer runtime.GetConnector().Close(host)
	if err != nil {
		return "", errors.Wrapf(err, "failed to connect to %s", host.GetAddress())
	}

	r := &connector.Runner{
		Conn: conn,
		//Debug: runtime.Arg.Debug,
		Host: host,
	}

	respMsg, err := r.SudoCmd(cmd, false)
	if err != nil {
		return "", err
	}
	return respMsg, nil
}
