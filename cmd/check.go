package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/zhilv666/linkchecker/pkg/log"
	"go.uber.org/zap"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检测网盘链接状态",
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("file")
		output, _ := cmd.Flags().GetString("output")
		server, _ := cmd.Flags().GetString("server")
		token, _ := cmd.Flags().GetString("token")
		proxy, _ := cmd.Flags().GetString("proxy")
		parallel, _ := cmd.Flags().GetInt("parallel")

		if parallel < 1 {
			parallel = 1
		}
		cpuNum := runtime.NumCPU()
		maxLimit := cpuNum * 2
		if parallel > maxLimit {
			fmt.Printf("⚠️  提示: 设定的并发数 (%d) 过高，已自动调整为系统建议上限: %d (CPU核数 * 2)\n", parallel, maxLimit)
			parallel = maxLimit
		}

		checker, err := NewAgentChecker(proxy, server, token, parallel)
		if err != nil {
			log.Fatal("failed to create agent", zap.Error(err))
		}
		if input != "" {
			fmt.Printf("📂 进入文件批处理模式 (Input: %s, Output: %s)\n", input, output)
			checker.RunFileMode(input, output)
		} else if len(args) > 0 {
			fmt.Println("💻 进入命令行模式")
			checker.RunArgsMode(args)
		} else {
			cmd.Help()
		}

	},
}

func init() {
	checkCmd.Flags().StringP("file", "f", "", "输入的 CSV 文件路径 (格式: 链接,密码)")
	checkCmd.Flags().StringP("output", "o", "./result.csv", "输出的 CSV 文件路径 (仅在文件模式下生效)")
	checkCmd.Flags().StringP("server", "s", "", "服务端上报接口地址 (e.g. http://127.0.0.1:3000)")
	checkCmd.Flags().StringP("token", "t", "", "服务端鉴权 Token")
	checkCmd.Flags().StringP("proxy", "x", "", "代理地址 (e.g. http://127.0.0.1:7890)")
	checkCmd.Flags().IntP("parallel", "p", 1, "并发检测数量")
	rootCmd.AddCommand(checkCmd)
}
