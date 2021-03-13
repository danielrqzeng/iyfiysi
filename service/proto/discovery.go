package proto

import (
	"iyfiysi/service"
	"strings"
	"time"
)

type GatewayParams struct {
	PackageName  string
	CreateTime   time.Time
	Domain       string
	AppName      string         `json:"app_name"`
	ServicesList []RpcParamType //要去重,此得到的是rpc名称
	RpcList      []RpcParamType
}

func genDiscoveryFile(
	tmplFile string,
	packageName, domain, appName string,
	rpcs []*RpcInfo) (buffStr string, err error) {
	params := &GatewayParams{}
	params.PackageName = packageName
	params.CreateTime = time.Now()
	params.Domain = domain
	params.AppName = appName
	params.RpcList = make([]RpcParamType, 0)
	params.ServicesList = make([]RpcParamType, 0)
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
		return
	}

	tmplStr, err := service.GetTmpl(tmplFile)
	if err != nil {
		return
	}
	buffStr, err = DoTmpl(tmplStr, params)
	return
}
