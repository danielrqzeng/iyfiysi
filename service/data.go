package service

import (
	"time"
)

const (
	//配置的类型
	TemplateConfigFlagNone    = 0
	TemplateConfigFlagFile    = 1
	TemplateConfigFlagDir     = 2
	TemplateConfigFlagIyfiysi = 4
	TemplateConfigFlagProtoc  = 8
)

//ProjectFileType 项目文件结构
type ProjectFileType struct {
	ID     string                 `json:"id"`     //id
	Flag   int                    `json:"flag"`   //类型，TemplateConfigType*,标识(比特组合)，0:none,1:文件类型,2:目录类型,4:给iyfiysi使用,8:给protoc-gen-iyfiysi使用
	Name   string                 `json:"name"`   //名称
	Desc   string                 `json:"desc"`   //描述
	Src    string                 `json:"src"`    //对应哪个template文件
	Dst    string                 `json:"dst"`    //生成之后放在那个文件
	Params map[string]interface{} `json:"params"` //参数，kv格式
}

//TemplateParams 模板参数
type TemplateParams struct {
	CreateTime time.Time //创建时间
	Domain     string    //域名
	AppName    string    //App名称
	Param1     string    //额外的
}
