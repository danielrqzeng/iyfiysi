package component

import (
	"path/filepath"
	"time"
)

func CreateGatewayScript(projectBase string) (err error) {
	CreateGatewayScriptStart(projectBase)
	CreateGatewayScriptStop(projectBase)
	CreateGatewayScriptCheck(projectBase)
	CreateGatewayScriptInclude(projectBase)
	return
}

func CreateGatewayScriptStart(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "gateway", "script", "start.sh")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("gateway_script_start.sh.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

func CreateGatewayScriptStop(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "gateway", "script", "stop.sh")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("gateway_script_stop.sh.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

func CreateGatewayScriptCheck(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "gateway", "script", "check.sh")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("gateway_script_check.sh.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

func CreateGatewayScriptInclude(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "gateway", "script", "include")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("gateway_script_include.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}
