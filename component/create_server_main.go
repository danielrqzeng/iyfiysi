package component

import (
	"path/filepath"
	"time"
)

func CreateServerMain(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "server", "main.go")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}

	tmplStr, err := GetTmpl("server_main.go.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption(DoFormat()))
	return
}
