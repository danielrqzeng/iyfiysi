package component

import (
	"path/filepath"
	"time"
)

func CreateGatewayDiscoveryInit(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "gateway", "discovery", "init.go")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("gateway_discovery_init.go.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption(DoFormat()))
	return
}

func CreateGatewayDiscoveryMain(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "gateway", "discovery", "main.go")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}
	tmplStr, err := GetTmpl("gateway_discovery_main.go.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption(DoFormat()))
	return
}

func GatewayDiscoveryDiscoveryNull(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "gateway", "discovery", "discovery.go")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    appName,
	}

	tmplStr, err := GetTmpl("gateway_discovery_discovery.go.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption(DoFormat()))
	return
}
