package iyfiysi

import (
	"github.com/spf13/cobra"
	"iyfiysi/internal/comm"
)

var (
	projectDomain = ""
	appName       = ""

	newCmd = &cobra.Command{
		Use:   "new",
		Short: "new a project",
		Long:  "new a project",
		Run: func(c *cobra.Command, args []string) {
			comm.Gen(projectDomain, appName)

		},
	}
)

/**
 * @brief 解析命令参数
 */
func init() {
	newCmd.PersistentFlags().StringVarP(&appName, "app", "a", "", "app name")
	newCmd.PersistentFlags().StringVarP(&projectDomain, "org", "o", "", "organization, a domain like 'iyfiysi.com'")
	//.\iyfiysi.exe new -o iyfiysi.com -a surl
}
