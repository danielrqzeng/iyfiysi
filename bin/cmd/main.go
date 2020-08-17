package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "iyfiysi",
		Short:        "iyfiysi tool",
		Long:         "iyfiysi tool",
		SilenceUsage: true,
	}
)

/**
 * @brief 初始化命令行工具
 */

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(protoCmd)
	rootCmd.AddCommand(versionCmd)
}

/**
 * @brief 执行命令行解析
 */
func Execute() {
	rootCmd.Execute()
}
