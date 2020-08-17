package cmd

import (
	"github.com/spf13/cobra"
	"iyfiysi/component"
)

var (
	projectDomain = ""
	projectName   = ""
	projectPath   = "."

	newCmd = &cobra.Command{
		Use:   "new",
		Short: "new a project",
		Long:  "new a project",
		Run: func(c *cobra.Command, args []string) {
			print("newCmd")
			component.CreateProject(projectDomain, projectName)
		},
	}
)

/**
 * @brief 解析命令参数
 */
func init() {
	newCmd.PersistentFlags().StringVarP(&projectPath, "dir", "d", ".", "project dir")
	newCmd.PersistentFlags().StringVarP(&projectName, "project", "p", "", "project name")
	newCmd.PersistentFlags().StringVarP(&projectDomain, "domain", "n", "", "project domain,go mod need this")
	//.\iyfiysi.exe new -n test.com -p surl
}
