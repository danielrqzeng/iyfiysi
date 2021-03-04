package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"iyfiysi/service"
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
			fmt.Printf("version: %v\ncommit: %v\ndate: %v\n", version, commit, date)
			service.Gen("test.com", "short_url")
		},
	}
)
