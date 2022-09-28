package cmd

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
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

			host := GetHostsByName(name)
			if host == nil {
				logger.Log.Warnf("未找到主机：%s", name)
				return
			}

			conn, err := runtime.GetConnector().Connect(host)
			defer runtime.GetConnector().Close(host)
			if err != nil {
				logger.Log.Errorln(errors.Wrapf(err, "failed to connect to %s", name))
				return
			}
			sess, err := conn.Session()
			if err != nil {
				logger.Log.Errorln(errors.Wrap(err, "failed to get SSH session"))
				return
			}
			defer sess.Close()

			exitCode := 0

			in, _ := sess.StdinPipe()
			out, _ := sess.StdoutPipe()

			err = sess.Shell()
			if err != nil {
				exitCode = -1
				if exitErr, ok := err.(*ssh.ExitError); ok {
					exitCode = exitErr.ExitStatus()
				}
				logger.Log.Errorln(err, exitCode)
				return
			}

			fmt.Println("连接主机...")
			fmt.Println("连接主机成功")
			fmt.Println("Last login:", time.Now())
			fmt.Printf("[root@%s ~]# ", host.GetName())

			reader := bufio.NewReader(os.Stdin) //终端输入
			var command string
			for {
				for {
					b, err := reader.ReadByte()
					if err != nil {
						break
					}
					if b == byte('\n') {
						break
					}
					command += string(b)
				}

				if command != "" {
					_, err = in.Write([]byte(command + "\n")) //主机输入
					if err != nil {
						logger.Log.Errorln(err)
						return
					}

					var (
						output []byte //存储主机输出
						line   = ""
						r      = bufio.NewReader(out) //主机输出
					)

					for {
						b, err := r.ReadByte() //主机输出
						if err != nil {
							break
						}

						output = append(output, b) //读取到的每个byte放到output中

						if b == byte('\n') {
							line = ""
							continue
						}

						line += string(b)

						if (strings.HasPrefix(line, "[sudo] password for ") || strings.HasPrefix(line, "Password")) && strings.HasSuffix(line, ": ") {
							_, err = in.Write([]byte(host.GetPassword() + "\n")) //主机输入
							if err != nil {
								break
							}
						}
					}

					err = sess.Wait()
					if err != nil {
						exitCode = -1
						if exitErr, ok := err.(*ssh.ExitError); ok {
							exitCode = exitErr.ExitStatus()
						}
					}
					outStr := strings.TrimPrefix(string(output), fmt.Sprintf("[sudo] password for %s:", host.GetUser()))

					logger.Log.Messagef(host.GetName(), outStr)
					//fmt.Println(command)
				}

				fmt.Printf("[%s@%s ~]# ", host.GetUser(), host.GetName())
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
