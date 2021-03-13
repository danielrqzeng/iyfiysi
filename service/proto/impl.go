package proto

import (
	"iyfiysi/component"
	"iyfiysi/service"
	"path/filepath"
	"strings"
	"time"
)

//生成服务实现代码
type ServerServiceRpcParams struct {
	//field for gen
	PackageName  string
	CreateTime   time.Time
	Domain       string
	AppName      string `json:"app_name"`
	ServiceName  string
	MethodName   string
	RequestName  string
	ResponseName string
}

func genServiceRpcFile(packageName string, domain, appName string, rpcs []*RpcInfo) {
	for _, rpc := range rpcs {
		params := &ServerServiceRpcParams{}
		params.PackageName = packageName
		params.Domain = domain
		params.AppName = appName
		params.CreateTime = time.Now()
		params.ServiceName = rpc.ServiceName
		params.MethodName = rpc.RpcName
		params.RequestName = rpc.RequestName
		params.ResponseName = rpc.ResponseName

		fileName := strings.ToLower(rpc.ServiceName) + "." + strings.ToLower(rpc.RpcName) + ".go"
		absFile := filepath.Join("..", "internal", "app", "server", "service", fileName)
		tmplStr, err := component.GetTmpl("server_service_service_rpc_protoc.go.tmpl")
		if err != nil {
			return
		}
		err = component.DoWriteFile(tmplStr, params, absFile, component.NewDoWriteFileOption(component.DoFormat()))
	}
}

func genServiceImplRpc(
	tmplFile string,
	packageName, domain, appName string,
	rpc *RpcInfo) (buffStr string, err error) {
	params := &ServerServiceRpcParams{}
	params.PackageName = packageName
	params.Domain = domain
	params.AppName = appName
	params.CreateTime = time.Now()
	params.ServiceName = rpc.ServiceName
	params.MethodName = rpc.RpcName
	params.RequestName = rpc.RequestName
	params.ResponseName = rpc.ResponseName

	tmplStr, err := service.GetTmpl(tmplFile)
	if err != nil {
		return
	}
	buffStr, err = DoTmpl(tmplStr, params)
	return
}
