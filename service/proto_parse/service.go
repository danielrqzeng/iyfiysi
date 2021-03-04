package proto_parse

import (
	"fmt"
	"iyfiysi/component"
	"path/filepath"
	"strings"
	"time"
)

type ServerParams struct {
	PackageName  string
	CreateTime   time.Time
	Domain       string
	AppName      string          `json:"app_name"`
	ServicesList []ServiceParams //要去重,此得到的是rpc名称
	RpcList      []ServiceParams //此是得到rpc的服务名称
}

func genServiceFile(projectBase string, domain, appName string, rpcs []*RpcInfo) {
	serverParams := &ServerParams{}
	serverParams.PackageName = "service"
	serverParams.CreateTime = time.Now()
	serverParams.Domain = domain
	serverParams.AppName = appName
	serverParams.ServicesList = make([]ServiceParams, 0)
	serverParams.RpcList = make([]ServiceParams, 0)
	serviceExist := make(map[string]bool)
	for _, rpc := range rpcs {
		for _, route := range rpc.Paths {
			params := ServiceParams{}
			params.ServiceName = strings.Title(rpc.ServiceName)
			params.Route = route
			params.MethodName = rpc.RpcName
			params.RequestName = rpc.RequestName
			params.ResponseName = rpc.ResponseName
			serverParams.RpcList = append(serverParams.RpcList, params)
			if _, ok := serviceExist[rpc.ServiceName]; !ok {
				serviceExist[rpc.ServiceName] = true
				serverParams.ServicesList = append(serverParams.ServicesList, params)
			}
		}
	}
	if len(serverParams.RpcList) == 0 {
		fmt.Println("not found services")
		return
	}

	absFile := filepath.Join("..", "internal", "app", "server", "service", "service.go")
	tmplStr, err := component.GetTmpl("server_service_service_protoc.go.tmpl")
	if err != nil {
		return
	}
	fmt.Println(tmplStr)
	fmt.Println("----------------------")
	err = component.DoWriteFile(tmplStr, serverParams, absFile, component.NewDoWriteFileOption(component.DoFormat()))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("done")
	}
}
