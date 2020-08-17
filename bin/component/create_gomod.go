package component

import (
	"path/filepath"
	"time"
)

func CreateGoMod(projectBase string) (err error) {
	output2temple := map[string]string{
		filepath.Join(projectBase, "go.mod"):            "project_go.mod.tmpl",
		filepath.Join(projectBase, "gateway", "go.mod"): "gateway_go.mod.tmpl",
		filepath.Join(projectBase, "server", "go.mod"):  "server_go.mod.tmpl",
	}

	for outputFile, tmplFile := range output2temple {
		domain, appName := GetDomainAppName(projectBase)
		c := &DefaultParams{
			CreateTime: time.Now(),
			Domain:     domain,
			AppName:    appName,
		}
		tmplStr := ""
		tmplStr, err = GetTmpl(tmplFile)
		if err != nil {
			return
		}
		err = DoWriteFile(tmplStr, c, outputFile, NewDoWriteFileOption())
	}
	return
}
