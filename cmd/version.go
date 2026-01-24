package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/zhilv666/linkchecker/internal/conf"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		goVersion := fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)

		fmt.Printf(`Build At: %s
Go Version: %s
Author: %s
Email: %s
Commit ID: %s
Version: %s
`, conf.BuildAt, goVersion, conf.GitAuthor, conf.GitEmail, conf.GitCommit, conf.Version)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
