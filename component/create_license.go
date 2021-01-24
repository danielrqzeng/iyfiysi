package component

import (
	"path/filepath"
	"strings"
	"time"
)

func CreateLicense(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "LICENSE")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    strings.Title(appName),
	}
	tmplStr, err := GetTmpl("LICENSE.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}
