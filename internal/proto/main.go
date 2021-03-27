package proto

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/iancoleman/strcase"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"io"
	"io/ioutil"
	"iyfiysi/internal/comm"
	"iyfiysi/util"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const usage = `
protoc -I. --iyfiysi_out=[OPTION,OPTION,...]:. ${pbfile.proto}
OPTION:
	domain=xxx.com
	app=test
e.g.
	protoc -I. --iyfiysi_out=logtostderr=true,domain=xxx.com,app=test:. ./service.proto
`

type RpcInfo struct {
	ServiceName  string
	RpcName      string
	RequestName  string
	ResponseName string
	Methods      []string
	Paths        []string
}

type RpcParamType struct {
	ServiceName  string
	Route        string
	MethodName   string
	RequestName  string
	ResponseName string
}

//pb文件的参数信息
type ProtoParamType struct {
	Services []string //要去重,此得到的是service list
	RpcList  []RpcParamType
}

func (this *RpcInfo) String() string {
	str := "\n"
	for i, method := range this.Methods {
		path_ := this.Paths[i]
		str += fmt.Sprintf(">  %s %s=>%s::%s(%s,%s)\n",
			method, path_, this.ServiceName, this.RpcName, this.RequestName, this.ResponseName)
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

func DoTmpl(tmplStr string, params interface{}) (targetStr string, err error) {
	//输出到buf
	buf := new(bytes.Buffer)
	// 创建模板对象, parse关联模板
	tmpl := template.Must(template.New(util.Md5sum([]byte(tmplStr))).Parse(tmplStr))
	err = tmpl.Execute(buf, params)
	if err != nil {
		return
	}
	targetStr = buf.String()
	return
}

func DoParse() {
	flag.Parse()
	req, err := ParseRequest(os.Stdin)
	if err != nil {
		glog.Fatal(err)
	}
	err = comm.InitGenConfig()
	if err != nil {
		glog.Fatal(err)
	}

	//glog.V(1).Info("Parsed code generator request")

	//glog.Error(req.GetParameter())

	//读取protoc传进来的参数,主要有两个,domain&app
	domain, appName := "", ""
	if req.Parameter != nil {
		for _, p := range strings.Split(req.GetParameter(), ",") {
			spec := strings.SplitN(p, "=", 2)
			if len(spec) == 1 {
				glog.Error(usage)
				os.Exit(1)
				return
			}
			name, value := spec[0], spec[1]
			switch name {
			case "domain":
				domain = value
			case "app":
				appName = value
			default:
				glog.Error(usage)
				os.Exit(1)
				return
			}
		}
	}

	//根据pb文件,解析出来协议信息
	rpcs := make([]*RpcInfo, 0)
	for _, protofile := range req.ProtoFile {
		for _, service := range protofile.Service {
			for _, rpc := range service.Method {
				//glog.V(1).Infof("method=%s,req=%,rsp=%s", *rpc.Name, rpc.GetInputType(), rpc.GetOutputType())
				opts, _ := extractAPIOptions(rpc)
				rpcInfo := &RpcInfo{}
				rpcInfo.ServiceName = service.GetName()
				rpcInfo.RpcName = rpc.GetName()
				requests := strings.Split(rpc.GetInputType(), ".")
				responses := strings.Split(rpc.GetOutputType(), ".")
				rpcInfo.RequestName = requests[len(requests)-1]
				rpcInfo.ResponseName = responses[len(responses)-1]
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

	//模板参数有来自三个地方
	//1.必须参数,来自命令输入,比如domain,appName,来自程序生成,比如createTs
	//2.来自pb文件,其一般指的是各种服务,rpc等
	//3.来自配置文件template/template.yaml,一般有包名
	//pb文件的参数
	pbParams := &ProtoParamType{}
	pbParams.Services = make([]string, 0)
	isServiceInSet := make(map[string]bool) //[serviceName]=true|false
	for _, r := range rpcs {
		rpc := RpcParamType{}
		rpc.ServiceName = generator.CamelCase(r.ServiceName)
		rpc.Route = strings.Join(r.Paths, ",")
		rpc.MethodName = generator.CamelCase(r.RpcName)
		rpc.RequestName = generator.CamelCase(r.RequestName)
		rpc.ResponseName = generator.CamelCase(r.ResponseName)
		pbParams.RpcList = append(pbParams.RpcList, rpc)
		if _, ok := isServiceInSet[r.ServiceName]; ok {
			continue
		}
		isServiceInSet[r.ServiceName] = true
		pbParams.Services = append(pbParams.Services, generator.CamelCase(r.ServiceName))
	}

	//默认自带值
	defaultParams := make(map[string]interface{})
	defaultParams["Domain"] = domain
	defaultParams["AppName"] = appName
	defaultParams["CreateTime"] = time.Now()

	//读取配置文件
	var tmplList []comm.ProjectFileType
	err = viper.UnmarshalKey("templates", &tmplList)
	if err != nil {
		glog.Error("failed to Unmarshal config,err=" + aurora.Red(err.Error()).String())
		return
	}
	for _, tmpl := range tmplList {
		//此配置不是给protoc_iyfiysi_out使用的,跳过
		if tmpl.Flag&comm.TemplateConfigFlagProtoc == 0 {
			continue
		}
		//服务发现
		if tmpl.ID == "protoc_discovery" {
			//整合三个地方的变量
			params := make(map[string]interface{})
			//1.来自pb的变量
			byteStr, err := json.Marshal(pbParams)
			if err != nil {
				glog.Error("failed to Marshal pb params,err=" + aurora.Red(err.Error()).String())
				return
			}
			mapPBParams, err := util.JsonStr2Map(string(byteStr))
			for k, v := range mapPBParams {
				params[k] = v
			}
			//2.来自配置文件的变量
			for k, v := range defaultParams {
				params[k] = v
			}
			//2.来自配置文件的变量
			for k, v := range tmpl.Params {
				params[k] = v
			}
			tmplStr, err := comm.GetTmpl(tmpl.Src)
			if err != nil {
				glog.Error(err, "failed to get template="+tmpl.Src)
				return
			}
			buffStr, err := DoTmpl(tmplStr, params)

			response := new(plugin.CodeGeneratorResponse)
			fname := filepath.FromSlash(tmpl.Dst)
			response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(fname),
				Content: proto.String(buffStr),
			})
			//glog.Error("to output to =" + tmpl.Dst)

			data, err := proto.Marshal(response)
			if err != nil {
				glog.Error(err, "failed to marshal output proto")
			}
			_, err = os.Stdout.Write(data)
			if err != nil {
				glog.Error(err, "failed to write output proto")
			}
			glog.Error("protoc-gen-iyfiysi generated success,file=" + aurora.Green(fname).String())
		}
		//服务注册
		if tmpl.ID == "protoc_register" {
			//整合三个地方的变量
			params := make(map[string]interface{})
			//1.来自pb的变量
			byteStr, err := json.Marshal(pbParams)
			if err != nil {
				glog.Error("failed to Marshal pb params,err=" + aurora.Red(err.Error()).String())
				return
			}
			mapPBParams, err := util.JsonStr2Map(string(byteStr))
			for k, v := range mapPBParams {
				params[k] = v
			}
			//2.来自配置文件的变量
			for k, v := range defaultParams {
				params[k] = v
			}
			//2.来自配置文件的变量
			for k, v := range tmpl.Params {
				params[k] = v
			}
			tmplStr, err := comm.GetTmpl(tmpl.Src)
			if err != nil {
				glog.Error(err, "failed to get template="+tmpl.Src)
				return
			}
			buffStr, err := DoTmpl(tmplStr, params)

			response := new(plugin.CodeGeneratorResponse)
			fname := filepath.FromSlash(tmpl.Dst)
			response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(fname),
				Content: proto.String(buffStr),
			})

			data, err := proto.Marshal(response)
			if err != nil {
				glog.Error(err, "failed to marshal output proto")
			}
			_, err = os.Stdout.Write(data)
			if err != nil {
				glog.Error(err, "failed to write output proto")
			}
			glog.Error("protoc-gen-iyfiysi generated success,file=" + aurora.Green(fname).String())
		}
		//服务实现
		if tmpl.ID == "protoc_impl" {
			//整合三个地方的变量
			params := make(map[string]interface{})
			//2.来自配置文件的变量
			for k, v := range defaultParams {
				params[k] = v
			}
			//2.来自配置文件的变量
			for k, v := range tmpl.Params {
				params[k] = v
			}
			tmplStr, err := comm.GetTmpl(tmpl.Src)
			if err != nil {
				glog.Error(err, "failed to get template="+tmpl.Src)
				return
			}
			for _, rpc := range pbParams.RpcList {
				//1.来自pb的变量
				byteStr, err := json.Marshal(rpc)
				if err != nil {
					glog.Error("failed to Marshal pb params,err=" + aurora.Red(err.Error()).String())
					return
				}
				mapPBParams, err := util.JsonStr2Map(string(byteStr))
				for k, v := range mapPBParams {
					params[k] = v
				}
				buffStr, err := DoTmpl(tmplStr, params)

				response := new(plugin.CodeGeneratorResponse)
				baseName := fmt.Sprintf("%s_%s.go", strcase.ToSnake(rpc.ServiceName), strcase.ToSnake(rpc.MethodName))
				fname := filepath.Join(filepath.FromSlash(tmpl.Dst), baseName)
				//若是已经存在文件，则略过，以免覆盖掉用户的实现逻辑
				if util.IsPathExist(fname) {
					glog.Error("protoc-gen-iyfiysi skip generated,file=" + aurora.Yellow(fname).String())
					continue
				}
				response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
					Name:    proto.String(fname),
					Content: proto.String(buffStr),
				})

				data, err := proto.Marshal(response)
				if err != nil {
					glog.Error(err, "failed to marshal output proto")
				}
				_, err = os.Stdout.Write(data)
				if err != nil {
					glog.Error(err, "failed to write output proto")
				}
				glog.Error("protoc-gen-iyfiysi generated success,file=" + aurora.Green(fname).String())
			}

		}
	}
}
