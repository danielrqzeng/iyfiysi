package cmd

import (
	"github.com/spf13/cobra"
	"iyfiysi/component"
)

var (
	protoCmd = &cobra.Command{
		Use:   "proto",
		Short: "gen proto",
		Long:  "gen proto",
		Run: func(c *cobra.Command, args []string) {
			print("protoCmd")
			component.ParsePB()
		},
	}
)

/**
 * @brief 解析命令参数
 */
func init() {
	protoCmd.PersistentFlags().StringVarP(&projectPath, "dir", "d", ".", "project dir")
	protoCmd.PersistentFlags().StringVarP(&projectName, "project", "p", "", "project name")
}
