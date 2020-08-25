package component

import (
	"path/filepath"
	"strings"
	"time"
)

func CreateProtoGenShell(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "proto", "gen.sh")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    strings.Title(appName),
	}

	tmplStr, err := GetTmpl("proto_gen.sh.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

func CreateProtoNull(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "proto", "service.proto")
	domain, appName := GetDomainAppName(projectBase)
	c := &DefaultParams{
		CreateTime: time.Now(),
		Domain:     domain,
		AppName:    strings.Title(appName),
	}

	tmplStr, err := GetTmpl("proto_service.proto.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

func CreateDependentProto(projectBase string) (err error) {
	output2temple := map[string]string{
		filepath.Join(projectBase, "proto", "google", "api", "annotations.proto"):     "proto_google_api_annotations.proto.tmpl",
		filepath.Join(projectBase, "proto", "google", "api", "http.proto"):            "proto_google_api_http.proto.tmpl",
		filepath.Join(projectBase, "proto", "google", "api", "httpbody.proto"):        "proto_google_api_httpbody.proto.tmpl",
		filepath.Join(projectBase, "proto", "google", "protobuf", "descriptor.proto"): "proto_google_protobuf_descriptor.proto.tmpl",
		filepath.Join(projectBase, "proto", "google", "rpc", "code.proto"):            "proto_google_rpc_code.proto.tmpl",
		filepath.Join(projectBase, "proto", "google", "rpc", "error_details.proto"):   "proto_google_rpc_error_details.proto.tmpl",
		filepath.Join(projectBase, "proto", "google", "rpc", "status.proto"):          "proto_google_rpc_status.proto.tmpl",
	}

	for outputFile, tmplFile := range output2temple {
		domain, appName := GetDomainAppName(projectBase)
		c := &DefaultParams{
			CreateTime: time.Now(),
			Domain:     domain,
			AppName:    strings.Title(appName),
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
