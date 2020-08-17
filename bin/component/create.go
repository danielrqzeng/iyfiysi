package component

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"iyfiysi/util"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func GetTmpl(tmplName string) (tmplStr string, err error) {
	tmplByte, err := ioutil.ReadFile(filepath.Join("template", tmplName))
	if err != nil {
		return
	}
	tmplStr = string(tmplByte)
	return
}

type DoWriteFileOption struct {
	DoFormat bool
}

type Option func(msg *DoWriteFileOption)

func DoFormat() Option {
	return func(m *DoWriteFileOption) {
		m.DoFormat = true
	}
}

func NewDoWriteFileOption(opts ...Option) *DoWriteFileOption {
	m := DoWriteFileOption{}
	for _, o := range opts {
		o(&m)
	}
	return &m
}

type DefaultParams struct {
	CreateTime time.Time
	Domain     string
	AppName    string
}

func DoWriteFile(tmplStr string, params interface{}, absFile string, opts *DoWriteFileOption) (err error) {
	//输出到目标文件
	targetWriter, err := os.OpenFile(absFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("open failed err:", err)
		return
	}

	// 创建模板对象, parse关联模板
	tmpl := template.Must(template.New(util.Md5sum([]byte(tmplStr))).Parse(tmplStr))
	err = tmpl.Execute(targetWriter, params)
	if err != nil {
		return
	}
	targetWriter.Close()

	if opts.DoFormat {
		//做一下格式化
		var beforeFormat, afterFormat []byte
		beforeFormat, err = ioutil.ReadFile(absFile)
		if err != nil {
			return
		}
		afterFormat, err = format.Source(beforeFormat)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ioutil.WriteFile(absFile, afterFormat, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	return
}

func GetDomainAppName(projectBase string) (domainName, appName string) {
	appName = filepath.Base(projectBase)
	domainName = filepath.Base(strings.TrimSuffix(projectBase, appName))
	return
}

func CreateProject(projectDomain, projectName string) {
	goBase := filepath.Join(util.GetGoPath(), "src")
	fmt.Println(goBase, projectName)
	util.DelPath(filepath.Join(goBase, projectDomain, projectName))
	//创建项目文件架构
	projectBase, err := CreateProjectPathStruct(projectName, filepath.Join(goBase, projectDomain, ))
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}
	//配置相关
	err = CreateConf(projectBase)
	err = CreateToolInit(projectBase)
	err = CreateToolJaeger(projectBase)
	err = CreateToolEtcdv3(projectBase)
	err = CreateUtilInit(projectBase)
	err = CreateUtilLogger(projectBase)
	err = CreateUtilUtil(projectBase)
	err = CreateKeystore(projectName, projectBase)
	//server
	err = CreateServerMain(projectBase)
	err = CreateServerServiceInit(projectBase)
	err = CreateServerServiceMain(projectBase)
	err = ServerServiceService(projectBase)
	//gateway
	err = CreateGatewayMain(projectBase)
	err = CreateGatewayDiscoveryInit(projectBase)
	err = CreateGatewayDiscoveryMain(projectBase)
	err = GatewayDiscoveryDiscoveryNull(projectBase)
	//proto
	err = CreateDependentProto(projectBase)
	err = CreateProtoNull(projectBase)
	err = CreateProtoGenShell(projectBase)
	//go.mod
	err = CreateGoMod(projectBase)
}
