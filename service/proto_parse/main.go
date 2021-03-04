package proto_parse

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
	"os"
	"strings"
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

func DoParse() {
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

	for _, r := range rpcs {
		fmt.Println(r.String())
	}

	domain, appName := "test.com", "short_url"
	genGatewayFile(`D:\go_path\src\surl`, domain, appName, rpcs)
	genServiceFile(`D:\go_path\src\surl`, domain, appName, rpcs)
	genServiceRpcFile(`D:\go_path\src\surl`, domain, appName, rpcs)
}
