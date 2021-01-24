package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"io"
	"io/ioutil"
	"iyfiysi/component"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RpcInfo struct {
	ServiceName  string
	RpcName      string
	RequestName  string
	ResponseName string
	Methods      []string
	Paths        []string
}

func (this *RpcInfo) String() string {
	str := "\n"
	for i, method := range this.Methods {
		path_ := this.Paths[i]
		str += fmt.Sprintf("%s %s=>%s::%s(%s,%s)\n", method, path_, this.ServiceName, this.RpcName, this.RequestName, this.ResponseName)
	}
	return str
}

// ParseRequest parses a code generator request from a proto Message.
func ParseRequest(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read code generator request: %v", err)
	}
	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(input, req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal code generator request: %v", err)
	}
	return req, nil
}

func extractAPIOptions(meth *descriptor.MethodDescriptorProto) (*options.HttpRule, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, options.E_Http) {
		return nil, nil
	}
	ext, err := proto.GetExtension(meth.Options, options.E_Http)
	if err != nil {
		return nil, err
	}
	opts, ok := ext.(*options.HttpRule)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
	}
	return opts, nil
}

func ParseOptionMethodPath(opt *options.HttpRule) (method, path_ string, err error) {
	method = ""
	path_ = ""
	if opt.GetPost() != "" {
		method = "POST"
		path_ = opt.GetPost()
	} else if opt.GetGet() != "" {
		method = "GET"
		path_ = opt.GetGet()
	} else if opt.GetPut() != "" {
		method = "PUT"
		path_ = opt.GetPut()
	} else if opt.GetDelete() != "" {
		method = "DELETE"
		path_ = opt.GetDelete()
	} else if opt.GetPatch() != "" {
		method = "PATCH"
		path_ = opt.GetPatch()
	}
	if method == "" {
		err = fmt.Errorf("not found method or path_")
		return
	}

	return
}

var (
	projectName = flag.String("project", "", "project name")
)

func main() {
	//flag.Set("log_dir", "/data/go_path/src") // 日志文件保存目录
	//flag.Set("v", "0")              // 配置V输出的等级
	//flag.Set("logtostderr", "true") // 配置V输出的等级
	flag.Parse()
	req, err := ParseRequest(os.Stdin)
	if err != nil {
		glog.Fatal(err)
	}
	glog.V(1).Info("Parsed code generator request")
	if req.Parameter != nil {
		for _, p := range strings.Split(req.GetParameter(), ",") {
			spec := strings.SplitN(p, "=", 2)
			if len(spec) == 1 {
				if err := flag.CommandLine.Set(spec[0], ""); err != nil {
					glog.Error(err)
					glog.Fatalf("Cannot set flag %s", p)
				}
				continue
			}
			name, value := spec[0], spec[1]
			glog.V(1).Info(name, value)
			if strings.HasPrefix(name, "M") {
				continue
			}
			if err := flag.CommandLine.Set(name, value); err != nil {
				glog.Fatalf("Cannot set flag %s", p)
			}
		}
	}

	rpcs := make([]*RpcInfo, 0)
	for _, protofile := range req.ProtoFile {
		for _, service := range protofile.Service {
			for _, rpc := range service.Method {
				//glog.V(1).Infof("method=%s,req=%,rsp=%s", *rpc.Name, rpc.GetInputType(), rpc.GetOutputType())
				opts, _ := extractAPIOptions(rpc)
				rpcInfo := &RpcInfo{}
				rpcInfo.ServiceName = service.GetName()
				rpcInfo.RpcName = rpc.GetName()
				rpcInfo.RequestName = rpc.GetInputType()[1:]
				rpcInfo.ResponseName = rpc.GetOutputType()[1:]
				//主方法
				mainMethod, mainPath, err := ParseOptionMethodPath(opts)
				if err != nil {
					panic(err)
				}

				rpcInfo.Methods = append(rpcInfo.Methods, mainMethod)
				rpcInfo.Paths = append(rpcInfo.Paths, mainPath)
				//附加方法
				extOpts := opts.GetAdditionalBindings()
				for _, extOpt := range extOpts {
					extMethod, extPath, err := ParseOptionMethodPath(extOpt)
					if err != nil {
						panic(err)
					}
					rpcInfo.Methods = append(rpcInfo.Methods, extMethod)
					rpcInfo.Paths = append(rpcInfo.Paths, extPath)

				}
				//glog.V(1).Info(rpcInfo.String())
				rpcs = append(rpcs, rpcInfo)
			}
		}
	}

	domain, appName := "test2.com", "surl"
	genGatewayFile(`D:\go_path\src\surl`, domain, appName, rpcs)
	genServiceFile(`D:\go_path\src\surl`, domain, appName, rpcs)
	genServiceRpcFile(`D:\go_path\src\surl`, domain, appName, rpcs)

}

type ServiceParams struct {
	//field for gen
	ServiceName  string
	Route        string
	MethodName   string
	RequestName  string
	ResponseName string
}

type GatewayParams struct {
	CreateTime   time.Time
	Domain       string
	AppName      string          `json:"app_name"`
	ServicesList []ServiceParams //要去重,此得到的是rpc名称
	RpcList      []ServiceParams
}

func genGatewayFile(projectBase string, domain, appName string, rpcs []*RpcInfo) {
	gatewayParams := &GatewayParams{}
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

	absFile := filepath.Join("..", "gateway", "discovery", "discovery.go")
	tmplStr, err := component.GetTmpl("gateway_discovery_discovery_protoc.go.tmpl")
	if err != nil {
		return
	}
	err = component.DoWriteFile(tmplStr, gatewayParams, absFile, component.NewDoWriteFileOption(component.DoFormat()))
}

type ServerServiceRpcParams struct {
	//field for gen
	CreateTime   time.Time
	Domain       string
	AppName      string `json:"app_name"`
	ServiceName  string
	MethodName   string
	RequestName  string
	ResponseName string
}

func genServiceRpcFile(projectBase string, domain, appName string, rpcs []*RpcInfo) {
	for _, rpc := range rpcs {
		params := &ServerServiceRpcParams{}
		params.Domain = domain
		params.AppName = appName
		params.CreateTime = time.Now()
		params.ServiceName = rpc.ServiceName
		params.MethodName = rpc.RpcName
		params.RequestName = rpc.RequestName
		params.ResponseName = rpc.ResponseName

		fileName := strings.ToLower(rpc.ServiceName) + "." + strings.ToLower(rpc.RpcName) + ".go"
		absFile := filepath.Join("..", "server", "service", fileName)
		tmplStr, err := component.GetTmpl("server_service_service_rpc_protoc.go.tmpl")
		if err != nil {
			return
		}
		err = component.DoWriteFile(tmplStr, params, absFile, component.NewDoWriteFileOption(component.DoFormat()))
	}
}

type ServerParams struct {
	CreateTime   time.Time
	Domain       string
	AppName      string          `json:"app_name"`
	ServicesList []ServiceParams //要去重,此得到的是rpc名称
	RpcList      []ServiceParams //此是得到rpc的服务名称
}

func genServiceFile(projectBase string, domain, appName string, rpcs []*RpcInfo) {
	serverParams := &ServerParams{}
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

	absFile := filepath.Join("..", "server", "service", "service.go")
	tmplStr, err := component.GetTmpl("server_service_service_protoc.go.tmpl")
	if err != nil {
		return
	}
	err = component.DoWriteFile(tmplStr, serverParams, absFile, component.NewDoWriteFileOption(component.DoFormat()))
}
