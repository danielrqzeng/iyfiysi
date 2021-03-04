package proto_parse

import (
	"fmt"
	"iyfiysi/component"
	"path/filepath"
	"strings"
	"time"
)

type ServiceParams struct {
	//field for gen
	ServiceName  string
	Route        string
	MethodName   string
	RequestName  string
	ResponseName string
}

type GatewayParams struct {
	PackageName  string
	CreateTime   time.Time
	Domain       string
	AppName      string          `json:"app_name"`
	ServicesList []ServiceParams //要去重,此得到的是rpc名称
	RpcList      []ServiceParams
}

func genGatewayFile(projectBase string, domain, appName string, rpcs []*RpcInfo) {
	gatewayParams := &GatewayParams{}
	gatewayParams.PackageName = "discovery"
	gatewayParams.CreateTime = time.Now()
	gatewayParams.Domain = domain
	gatewayParams.AppName = appName
	gatewayParams.RpcList = make([]ServiceParams, 0)
	gatewayParams.ServicesList = make([]ServiceParams, 0)
	serviceExist := make(map[string]bool)
	for _, rpc := range rpcs {
		for _, route := range rpc.Paths {
			params := ServiceParams{}
			params.ServiceName = strings.Title(rpc.ServiceName)
			params.Route = route
			params.MethodName = rpc.RpcName
			params.RequestName = rpc.RequestName
			params.ResponseName = rpc.ResponseName
			gatewayParams.RpcList = append(gatewayParams.RpcList, params)
			if _, ok := serviceExist[rpc.ServiceName]; !ok {
				serviceExist[rpc.ServiceName] = true
				gatewayParams.ServicesList = append(gatewayParams.ServicesList, params)
			}
		}
	}
	if len(gatewayParams.RpcList) == 0 {
		return
	}

	fmt.Println(gatewayParams)
	fmt.Println("----------------------")

	absFile := filepath.Join("..", "internal", "app", "gateway", "discovery", "discovery.go")
	tmplStr, err := component.GetTmpl("gateway_discovery_discovery_protoc.go.tmpl")
	if err != nil {
		return
	}
	fmt.Println(tmplStr)
	fmt.Println("----------------------")
	err = component.DoWriteFile(tmplStr, gatewayParams, absFile, component.NewDoWriteFileOption(component.DoFormat()))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("done")
	}
}
