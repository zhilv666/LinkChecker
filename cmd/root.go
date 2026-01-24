package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "linkchecker",
	Short: "LinkChecker 是一个用来检查网盘分享链接有效性和状态的工具。",
	Long:  `LinkChecker 是一个高效的工具，用于检测各大网盘分享链接的有效性、过期状态以及是否被封禁或删除。支持对百度网盘、夸克网盘多个网盘平台的链接进行分析，并提供详细的分享信息（如文件大小、创建者、过期时间等）。支持缓存机制，减少重复请求的负担，提升检测效率。该工具还支持命令行界面（CLI），方便开发者和用户快速集成和使用。`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
