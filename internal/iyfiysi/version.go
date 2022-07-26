package iyfiysi

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// Variables set at build time
var (
	version = "v1.0.0"
	commit  = "unknown"
	date    = "unknown"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version",
		Long:  "print version",
		Run: func(c *cobra.Command, args []string) {
			fmt.Printf("%s\t:%s\n%s\t:%s\n%s\t:%s\n",
				aurora.Green("version"), version,
				aurora.Green("commit"), commit,
				aurora.Green("date"), date)
			//fmt.Printf("version: %v\ncommit: %v\ndate: %v\n", version, commit, date)
		},
	}
)
