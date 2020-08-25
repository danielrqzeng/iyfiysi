package component

import (
	"path/filepath"
	"time"
)

func CreateToolInit(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "tool", "init.go")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}

	tmplStr, err := GetTmpl("tool_init.go.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption(DoFormat()))
	return
}

func CreateToolJaeger(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "tool", "jaeger.go")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("tool_jaeger.go.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption(DoFormat()))
	return
}

func CreateToolEtcdv3(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "tool", "etcdv3.go")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}

	tmplStr, err := GetTmpl("tool_etcdv3.go.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption(DoFormat()))
	return
}
