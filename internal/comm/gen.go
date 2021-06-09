package comm

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/viper"
	"go/format"
	"io/ioutil"
	"iyfiysi/util"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//GetTmpl 获取模板内容
//tmplName 是一个模板文件名（不包含目录）
func GetTmpl(tmplName string) (tmplStr string, err error) {
	/*	tmplByte, err := ioutil.ReadFile(filepath.Join("template", tmplName))
		if err != nil {
			return
		}
		tmplStr = string(tmplByte)
		return*/

	//使用statik打包下模板，以免生成bin文件后有依赖
	statikFS, err := fs.New()
	if err != nil {
		return
	}

	r, err := statikFS.Open("/" + tmplName)
	if err != nil {
		fmt.Println(tmplName, err)
		return
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	tmplStr = string(contents)
	return
}

//DoWriteFileOption 写入文件时候，带的参数
type DoWriteFileOption struct {
	DoFormat bool   //是否做格式化，true:yes
	Delims   string //模板分隔符，如果是空，则默认为{{}}
}

//Option 选项内容
type Option func(msg *DoWriteFileOption)

//DoFormat 选项-格式化
func DoFormat() Option {
	return func(m *DoWriteFileOption) {
		m.DoFormat = true
	}
}

//DoDelims 选项-分隔符
func DoDelims(delim string) Option {
	return func(m *DoWriteFileOption) {
		m.Delims = delim
	}
}

//NewDoWriteFileOption 新建一个DoWriteFileOption格式的option
func NewDoWriteFileOption(opts ...Option) *DoWriteFileOption {
	m := DoWriteFileOption{}
	for _, o := range opts {
		o(&m)
	}
	return &m
}

//DoWriteFile 生成模板并且写入文件
func DoWriteFile(tmplStr string, params interface{}, absFile string, opts *DoWriteFileOption) (err error) {
	//输出到目标文件
	targetWriter, err := os.OpenFile(absFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic("open file=" + absFile + " failed err:" + err.Error())
		return
	}
	// 创建模板对象, parse关联模板
	tmpl := template.New(util.Md5sum([]byte(tmplStr)))
	if opts.Delims != "" {
		delims := strings.Split(opts.Delims, ",")
		if len(delims) != 2 {
			panic("id=" + absFile + " err for delims=" + opts.Delims)
		}
		tmpl.Delims(delims[0], delims[1])
	}
	tmpl = template.Must(tmpl.Parse(tmplStr))
	err = tmpl.Execute(targetWriter, params)
	if err != nil {
		return
	}
	err = targetWriter.Close()
	if err != nil {
		return
	}

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

//InitGenConfig 初始化生成配置，其实主要是将template/template.yaml的配置读进viper里面
func InitGenConfig() (err error) {
	//获取配置
	const templateConfigFile = "template.yaml"
	configStr, err := GetTmpl(templateConfigFile)
	if err != nil {
		return
	}

	//读取&解析配置到viper中
	templateReader := strings.NewReader(configStr)
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(templateReader)
	if err != nil {
		fmt.Println("parse template/template.yaml fail ")
		return
	}
	return
	//viper.SetConfigFile("template/template.yaml")
	//err = viper.ReadInConfig()
}

func Mkdir(absDir string) (err error) {
	//先检查所有的目录都是存在的，若是不存在则新建
	if !util.IsPathExist(absDir) {
		err = os.MkdirAll(absDir, os.ModePerm)
		if err != nil {
			return
		}
		fmt.Println("mkdir baseDir ", aurora.Green(absDir), " success")
	} else {
		//fmt.Println("mkdir fail cuz baseDir ", aurora.Red(absDir), " exist")
	}
	return
}

//Gen 根据配置生成项目文件
//projectDir 项目的目录，格式为:[pre/dir/path/]domain.com,例如：/a/b/c/google.com/
//projectName 项目名称，格式为字符串
func Gen(projectDir, projectName string) (err error) {
	au := aurora.NewAurora(true)
	//读取配置文件template/template.yaml
	err = InitGenConfig()
	if err != nil {
		panic(err)
	}

	//检查projectDir是否合法,并且找到项目base目录
	baseDir := projectDir
	if !filepath.IsAbs(projectDir) {
		baseDir, err = filepath.Abs(projectDir)
		if err != nil {
			return
		}
	}

	baseDir = filepath.Join(baseDir, projectName)
	err = Mkdir(baseDir)
	if err != nil {
		panic(err)
	}

	var tmplList []ProjectFileType
	err = viper.UnmarshalKey("templates", &tmplList)
	if err != nil {
		panic(err)
	}
	for _, tmpl := range tmplList {
		//此配置不是给iyfiysi使用的,跳过
		if tmpl.Flag&TemplateConfigFlagIyfiysi == 0 {
			continue
		}
		//一个文件类型
		if tmpl.Flag&TemplateConfigFlagFile != 0 {
			err = GenFile(baseDir, &tmpl)
			if err != nil {
				//fmt.Println("id=" + tmpl.ID + ",err=" + err.Error())
				fmt.Println("id=", au.Red(tmpl.ID), ",err="+err.Error())
			} else {
				fmt.Println("success gen file ", au.Green(tmpl.ID))
			}
		}
		//一个目录类型
	}

	//生成证书
	c := viper.GetString("keystore.country")
	o := viper.GetString("keystore.organization")
	ou := viper.GetString("keystore.organizationalUnit")
	cn := viper.GetString("keystore.commonName")
	cacrt := viper.GetString("keystore.cacrt")
	cakey := viper.GetString("keystore.cakey")
	csr := viper.GetString("keystore.csr")
	crt := viper.GetString("keystore.crt")
	key := viper.GetString("keystore.key")
	expireDays := viper.GetInt("keystore.expireDays")
	dnsName := viper.GetStringSlice("keystore.dnsName")

	err = GenCert(baseDir,
		cacrt, cakey, csr, crt, key,
		c, o, ou, cn,
		dnsName,
		expireDays)
	if err != nil {
		fmt.Println("cert fail,", au.Red(err.Error()))
	}
	fmt.Println("success gen file ", au.Green(cacrt))
	fmt.Println("success gen file ", au.Green(cakey))
	fmt.Println("success gen file ", au.Green(csr))
	fmt.Println("success gen file ", au.Green(crt))
	fmt.Println("success gen file ", au.Green(key))
	return
}

//GenFile 根据模板，生成文件
//baseDir: 举例：/data/go_path/src/github.com/app
//templateFile: init.go.tmpl
//dstFile: internal/pkg/utils/init.go
func GenFile(baseDir string, fileConfig *ProjectFileType) (err error) {
	//生成目录
	dstFile := filepath.FromSlash(fileConfig.Dst)
	absFile := filepath.Join(baseDir, dstFile)
	absDir := filepath.Dir(absFile)
	err = Mkdir(absDir)
	if err != nil {
		return
	}

	//如果文件已經存在，则不能生成
	if util.IsPathExist(absFile) {
		err = fmt.Errorf("file=" + absFile + " exist")
		return
	}

	//直接复制
	if fileConfig.Flag&TemplateConfigFlagCopy != 0 {
		fileStr := ""
		fileStr, err = GetTmpl(fileConfig.Src)
		if err != nil {
			return
		}
		err = util.WriteFile(absFile, []byte(fileStr))
		if err != nil {
			return
		}
		return
	}
	//模板复制

	params := make(map[string]interface{})
	if fileConfig.Params != nil {
		params = fileConfig.Params
	}
	domainStr, appName := filepath.Split(baseDir)
	domainStr = filepath.Base(domainStr)
	params["AppName"] = appName //filepath.Base(baseDir)
	params["Domain"] = domainStr
	params["CreateTime"] = time.Now()

	tmplStr, err := GetTmpl(fileConfig.Src)
	if err != nil {
		return
	}

	opts := NewDoWriteFileOption(DoDelims(fileConfig.Delims))
	if strings.HasSuffix(absFile, ".go") && fileConfig.Delims != "" {
		opts = NewDoWriteFileOption(DoFormat(), DoDelims(fileConfig.Delims))
	}
	err = DoWriteFile(tmplStr, params, absFile, opts)
	if err != nil {
		return
	}

	return
}

//GenCert 生成证书
func GenCert(
	baseDir string,
	caPubFile, caKeyFile, csrFile, crtFile, keyFile string,
	country, organization, organizationalUnit, commonName string,
	dnsName []string,
	expireDays int) (err error) {

	//检查证书文件是否齐全
	files := []string{caPubFile, caKeyFile, csrFile, crtFile, keyFile}
	for i, fp := range files {
		//生成目录
		dstFile := filepath.FromSlash(fp)
		absFile := filepath.Join(baseDir, dstFile)
		absDir := filepath.Dir(absFile)
		err = Mkdir(absDir)
		if err != nil {
			return
		}
		//如果文件已經存在，则不能生成
		if util.IsPathExist(absFile) {
			err = fmt.Errorf("file=" + absFile + " exist,fail to gen new cert")
			return
		}
		switch i {
		case 0:
			caPubFile = absFile
		case 1:
			caKeyFile = absFile
		case 2:
			csrFile = absFile
		case 3:
			crtFile = absFile
		case 4:
			keyFile = absFile
		}
	}

	//生成ca
	err = CreateCA("CN", "IYFIYSI", "SECURE", "CERT",
		caPubFile, caKeyFile)
	if err != nil {
		return
	}
	//生成csr和密钥
	err = CreateCSR(country, organization, organizationalUnit, commonName,
		csrFile, keyFile)
	if err != nil {
		return
	}

	//生成crt
	err = CreateCert(country, organization, organizationalUnit, commonName,
		caPubFile, caKeyFile, csrFile, crtFile, keyFile,
		dnsName,
		expireDays)
	if err != nil {
		return
	}
	return
}
