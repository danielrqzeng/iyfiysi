package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"io"
	"io/ioutil"
	"iyfiysi/component"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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
	flag.Set("v", "0")              // 配置V输出的等级
	flag.Set("logtostderr", "true") // 配置V输出的等级
	flag.Parse()
	glog.V(1).Info("Parsing code generator request")
	req, err := ParseRequest(os.Stdin)
	if err != nil {
		glog.Fatal(err)
	}
	glog.V(1).Info("Parsed code generator request")
	fmt.Println("111 before parse req.Parameter")
	if req.Parameter != nil {
		fmt.Println("req.Parameter is not null")
		for _, p := range strings.Split(req.GetParameter(), ",") {
			glog.V(1).Info(p)
			fmt.Println("p=" + p)
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
	for _, target := range req.FileToGenerate {
		glog.V(1).Info(target)
	}
	rpcs := make([]*RpcInfo, 0)
	for _, protofile := range req.ProtoFile {
		glog.V(1).Info(*protofile.Name)
		for _, service := range protofile.Service {
			glog.V(1).Info(*service.Name)
			for _, rpc := range service.Method {
				glog.V(1).Infof("method=%s,req=%,rsp=%s", *rpc.Name, rpc.GetInputType(), rpc.GetOutputType())
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
				glog.V(1).Info(rpcInfo.String())
				rpcs = append(rpcs, rpcInfo)
				fmt.Println("---------")
				fmt.Println(rpcInfo.String())
			}
		}
	}

	domain, appName := "test1.com", "surl"
	genGatewayFile(`D:\go_path\src\surl`, domain, appName, rpcs)
	genServiceFile(`D:\go_path\src\surl`, domain, appName, rpcs)
	genServiceRpcFile(`D:\go_path\src\surl`, domain, appName, rpcs)

}

const gatewayTmpl = `// gen by iyfiysi at {{.CreateTime}}
package discovery

import (
	"context"
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "{{.Domain}}/{{.AppName}}/proto"
)

func DoDiscovery(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) (err error) {
    {{- range $Service := .Services}}
	//for {{$Service.Route}}=>{{$Service.ServiceName}}::{{$Service.MethodName}}({{$Service.RequestName}},{{$Service.ResponseName}})
    err = gw.Register{{$Service.ServiceName}}HandlerFromEndpoint(ctx, mux, "{{$Service.Route}}", opts)
    if err != nil {
        return
    }
    {{end}}
	return
}
`

type ServiceParams struct {
	//field for gen
	ServiceName  string
	Route        string
	MethodName   string
	RequestName  string
	ResponseName string
}

type GatewayParams struct {
	CreateTime time.Time
	Domain     string
	AppName    string `json:"app_name"`
	Services   []ServiceParams
}

func genGatewayFile(projectBase string, domain, appName string, rpcs []*RpcInfo) {
	gatewayParams := &GatewayParams{}
	gatewayParams.CreateTime = time.Now()
	gatewayParams.Domain = domain
	gatewayParams.AppName = appName
	gatewayParams.Services = make([]ServiceParams, 0)
	for _, rpc := range rpcs {
		for _, route := range rpc.Paths {
			params := ServiceParams{}
			params.ServiceName = strings.Title(rpc.ServiceName)
			params.Route = route
			params.MethodName = rpc.RpcName
			params.RequestName = rpc.RequestName
			params.ResponseName = rpc.ResponseName
			gatewayParams.Services = append(gatewayParams.Services, params)
		}
	}
	if len(gatewayParams.Services) == 0 {
		return
	}

	//targetWriter, err := os.OpenFile(filepath.Join("..", "gateway", "discovery", "discovery.go"), os.O_CREATE|os.O_WRONLY, 0755)
	//if err != nil {
	//	fmt.Println("open failed err:", err)
	//	return
	//}
	// 创建模板对象, parse关联模板
	//tmpl := template.Must(template.New("genGatewayFile").Parse(gatewayTmpl))
	//err := tmpl.Execute(targetWriter, gatewayParams)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	absFile := filepath.Join("..", "gateway", "discovery", "discovery.go")
	tmplStr, err := component.GetTmpl("gateway_discovery_discovery_protoc.go.tmpl")
	if err != nil {
		return
	}
	err = component.DoWriteFile(tmplStr, gatewayParams, absFile, component.NewDoWriteFileOption(component.DoFormat()))
}

const serverServiceRpcTmpl = `// gen by iyfiysi at {{.CreateTime}}
package service

import(
	"context"
	"fmt"
	"{{.Domain}}/{{.AppName}}/proto"
)

func (s *{{.ServiceName}}Impl) {{.MethodName}}(ctx context.Context, req *{{.RequestName}})  (rsp *{{.ResponseName}}, err error) {
	rsp =&{{.ResponseName}}{}
	fmt.Println(req)
	return
}
`

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

		fileName := strings.ToLower(rpc.ServiceName) + "." + strings.ToLower(rpc.RpcName) + ".go"
		targetWriter, err := os.OpenFile(filepath.Join("..", "server", "service", fileName), os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Println("open failed err:", err)
			return
		}
		params := &ServerServiceRpcParams{}
		params.Domain = domain
		params.AppName = appName
		params.CreateTime = time.Now()
		params.ServiceName = rpc.ServiceName
		params.MethodName = rpc.RpcName
		params.RequestName = rpc.RequestName
		params.ResponseName = rpc.ResponseName

		// 创建模板对象, parse关联模板
		tmpl := template.Must(template.New("genServiceRpcFile" + rpc.ServiceName + rpc.RpcName).Parse(serverServiceRpcTmpl))
		err = tmpl.Execute(targetWriter, params)
		if err != nil {
			return
		}
	}
}

const serverServiceTmpl = `// gen by iyfiysi at {{.CreateTime}}
package service

import (
    "google.golang.org/grpc"
    "github.com/spf13/viper"
    "{{.Domain}}/{{.AppName}}/tool"
    "{{.Domain}}/{{.AppName}}/proto"
)

{{- range $ServicesStruct := .ServicesStructs}}
type {{$ServicesStruct.ServiceName}}Impl struct{}
{{end}}

func DoRegister(grpcServer *grpc.Server) (err error) {
{{range $ServicesStruct := .ServicesStructs}}
	{
		s := &{{$ServicesStruct.ServiceName}}Impl{}
    	proto.Register{{$ServicesStruct.ServiceName}}Server(grpcServer, s)
	}
{{end}}

    instance := viper.GetString("server.listen")
{{- range $Service := .Services}}
    tool.Register("{{$Service.Route}}",instance)
{{end}}
    return
}
`

type ServerParams struct {
	CreateTime      time.Time
	Domain          string
	AppName         string          `json:"app_name"`
	ServicesStructs []ServiceParams //要去重
	Services        []ServiceParams
}

func genServiceFile(projectBase string, domain, appName string, rpcs []*RpcInfo) {
	fmt.Println("genServiceFile")
	serverParams := &ServerParams{}
	serverParams.CreateTime = time.Now()
	serverParams.Domain = domain
	serverParams.AppName = appName
	serverParams.ServicesStructs = make([]ServiceParams, 0)
	serverParams.Services = make([]ServiceParams, 0)
	serviceExist := make(map[string]bool)
	for _, rpc := range rpcs {
		for _, route := range rpc.Paths {
			params := ServiceParams{}
			params.ServiceName = strings.Title(rpc.ServiceName)
			params.Route = route
			params.MethodName = rpc.RpcName
			params.RequestName = rpc.RequestName
			params.ResponseName = rpc.ResponseName
			serverParams.Services = append(serverParams.Services, params)
			if _, ok := serviceExist[rpc.ServiceName]; !ok {
				serviceExist[rpc.ServiceName] = true
				serverParams.ServicesStructs = append(serverParams.ServicesStructs, params)
			}
		}
	}
	if len(serverParams.Services) == 0 {
		fmt.Println("not found services")
		return
	}

	targetWriter, err := os.OpenFile(filepath.Join("..", "server", "service", "service.go"), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("open failed err:", err)
		return
	}
	fmt.Println("save to " + filepath.Join("..", "server", "service", "service.go"))
	// 创建模板对象, parse关联模板
	tmpl := template.Must(template.New("genServiceFile").Parse(serverServiceTmpl))
	err = tmpl.Execute(targetWriter, serverParams)
	if err != nil {
		fmt.Println(err)
		return
	}
}
