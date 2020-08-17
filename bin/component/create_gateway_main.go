package component

import (
	"path/filepath"
	"time"
)

func CreateGatewayMain(projectBase string) (err error) {
	absFile := filepath.Join(projectBase,"gateway", "main.go")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}

	tmplStr, err := GetTmpl("gateway_main.go.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption(DoFormat()))
	return
}
