package component

import (
	"path/filepath"
	"time"
)

func CreateServerScript(projectBase string) (err error) {
	CreateServerScriptStart(projectBase)
	CreateServerScriptStop(projectBase)
	CreateServerScriptCheck(projectBase)
	CreateServerScriptInclude(projectBase)
	return
}

func CreateServerScriptStart(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "server", "script", "start.sh")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("server_script_start.sh.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

func CreateServerScriptStop(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "server", "script", "stop.sh")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("server_script_stop.sh.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

func CreateServerScriptCheck(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "server", "script", "check.sh")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("server_script_check.sh.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

func CreateServerScriptInclude(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "server", "script", "include")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("server_script_include.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}
