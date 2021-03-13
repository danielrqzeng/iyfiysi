package proto

import (
	"fmt"
	"iyfiysi/component"
	"iyfiysi/service"
	"iyfiysi/service/proto_parse"
	"path/filepath"
	"strings"
	"time"
)

//生成server的文件,用于做服务注册
type ServerParams struct {
	PackageName  string
	CreateTime   time.Time
	Domain       string
	AppName      string         `json:"app_name"`
	ServicesList []RpcParamType //要去重,此得到的是rpc名称
	RpcList      []RpcParamType //此是得到rpc的服务名称
}

func genServiceFile(packageName string, domain, appName string, rpcs []*proto_parse.RpcInfo) {
	serverParams := &ServerParams{}
	serverParams.PackageName = packageName
	serverParams.CreateTime = time.Now()
	serverParams.Domain = domain
	serverParams.AppName = appName
	serverParams.ServicesList = make([]RpcParamType, 0)
	serverParams.RpcList = make([]RpcParamType, 0)
	serviceExist := make(map[string]bool)
	for _, rpc := range rpcs {
		for _, route := range rpc.Paths {
			params := RpcParamType{}
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

func genRegisterFile(
	tmplFile string,
	packageName, domain, appName string,
	rpcs []*RpcInfo) (buffStr string, err error) {

	params := &ServerParams{}
	params.PackageName = packageName
	params.CreateTime = time.Now()
	params.Domain = domain
	params.AppName = appName
	params.ServicesList = make([]RpcParamType, 0)
	params.RpcList = make([]RpcParamType, 0)
	serviceExist := make(map[string]bool)
	for _, rpc := range rpcs {
		for _, route := range rpc.Paths {
			serviceParams := RpcParamType{}
			serviceParams.ServiceName = strings.Title(rpc.ServiceName)
			serviceParams.Route = route
			serviceParams.MethodName = rpc.RpcName
			serviceParams.RequestName = rpc.RequestName
			serviceParams.ResponseName = rpc.ResponseName
			params.RpcList = append(params.RpcList, serviceParams)
			if _, ok := serviceExist[rpc.ServiceName]; !ok {
				serviceExist[rpc.ServiceName] = true
				params.ServicesList = append(params.ServicesList, serviceParams)
			}
		}
	}
	if len(params.RpcList) == 0 {
		fmt.Println("not found services")
		return
	}

	tmplStr, err := service.GetTmpl(tmplFile)
	if err != nil {
		return
	}
	buffStr, err = DoTmpl(tmplStr, params)
	return
}
