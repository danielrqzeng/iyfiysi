package component

import (
	"fmt"
	"iyfiysi/util"
	"os"
	"path/filepath"
)

/*创建目录
projectName:项目名称，projectPath：项目路径
举例,假设gopath="/data/go_path",currPath="/data/go_path/src/iyfiysi"
1.projectName=project1,projectPath="/data/www",则输出为/data/www/project1,其使用绝对路径
2.projectName=project1,projectPath="../www",则输出为/data/go_path/src/www/project1，其使用了相对路径，且相对于程序运行的路径
1.projectName=project1,projectPath="",则输出为/data/go_path/src/project1，其使用的是gopath的路径
*/
func CreateDir(projectName string, projectPath string) (absPath string, err error) {
	//若是projectPath为空，则使用gopath
	if projectPath == "" {
		projectPath = util.GetGoPath()
		projectPath = filepath.Join(projectPath, "src")
	} else {
		//为相对路径，则相对于当前程序运行路径
		if !filepath.IsAbs(projectPath) {
			currPath := util.GetCurrPath()
			projectPath = filepath.Join(currPath, projectPath)
		}
	}
	targetPath := filepath.Join(projectPath, projectName)

	//若是存在，则报错
	if util.IsPathExist(targetPath) {
		err = fmt.Errorf("path=" + targetPath + " already exist")
		return
	}
	absPath = targetPath

	//os.MkdirAll(targetPath,os.ModePerm)
	err = os.MkdirAll(targetPath, 0755)
	return
}

func CreateProjectPathStruct(projectName string, projectPath string) (projectBase string, err error) {
	projectBase, err = CreateDir(projectName, projectPath)
	if err != nil {
		return
	}
	CreateDir("conf", projectBase)
	CreateDir("sql", projectBase)
	CreateDir("test", projectBase)
	CreateDir("keystore", projectBase)
	CreateDir("util", projectBase)
	CreateDir("tool", projectBase)
	CreateDir("data", projectBase)
	CreateDir("proto/google", projectBase)
	CreateDir("proto/google/api", projectBase)
	CreateDir("proto/google/protobuf", projectBase)
	CreateDir("proto/google/rpc", projectBase)
	CreateDir("server", projectBase)
	CreateDir("server/conf", projectBase)
	CreateDir("server/service", projectBase)
	CreateDir("gateway", projectBase)
	CreateDir("gateway/conf", projectBase)
	CreateDir("gateway/discovery", projectBase)
	return
}
