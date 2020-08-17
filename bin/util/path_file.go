package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

//获取运行程序的当前路径
func GetCurrPath() (currPath string) {
	currPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return
}

//获取GOPATH的路径
func GetGoPath() (goPath string) {
	goPath, ok := os.LookupEnv("GOPATH")
	if !ok {
		panic("env GOPATH not found")
	}
	return
}

func IsPathExist(pathName string) (exist bool) {
	exist = false
	_, err := os.Stat(pathName)
	if err == nil {
		exist = true
		return
	}
	return
}

func DelPath(pathName string) {
	err :=os.RemoveAll(pathName)
	if err != nil{
		fmt.Println(err)
	}
}

func CopyFile(srcFileName, dsrFileName string) (err error) {
	input, err := ioutil.ReadFile(srcFileName)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(dsrFileName, input, 0644)
	if err != nil {
		return
	}
	return
}
