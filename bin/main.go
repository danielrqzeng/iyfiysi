package main

import (
	"fmt"
	"golang.org/x/net/context"
	_ "google.golang.org/grpc"
	"iyfiysi/cmd"
	"iyfiysi/component"
	"iyfiysi/pb"
	"iyfiysi/util"
	_ "log"
	_ "net"
	"path/filepath"
)

// server is used to implement customer.CustomerServer.
type server struct {
}

// CreateCustomer creates a new Customer
func (s *server) BKService(ctx context.Context, in *pb.BlockRequest) (response *pb.BlockResponse, err error) {
	response = &pb.BlockResponse{}
	response.Message = in.Name
	fmt.Println(in)
	return
}

func CreateProject(projectName string) {
	goBase := filepath.Join(util.GetGoPath(), "src")
	util.DelPath(filepath.Join(goBase, projectName))
	//创建项目文件架构
	projectBase, err := component.CreateProjectPathStruct(projectName, "")
	if err != nil {
		fmt.Println(err)
	}
	//配置相关
	err = component.CreateConf(projectBase)
	err = component.CreateToolJaeger(projectBase)
	err = component.CreateUtilLogger(projectBase)
	err = component.CreateUtilUtil(projectBase)
	err = component.CreateKeystore(projectName, projectBase)
	//server
	err = component.CreateServerMain(projectBase)
	err = component.CreateServerServiceMain(projectBase)
	err = component.ServerServiceService(projectBase)
	//gateway
	err = component.CreateGatewayMain(projectBase)
	err = component.CreateGatewayDiscoveryMain(projectBase)
	err = component.GatewayDiscoveryDiscoveryNull(projectBase)
	//proto
	err = component.CreateDependentProto(projectBase)
	err = component.CreateProtoNull(projectBase)
	err = component.CreateProtoGenShell(projectBase)
}

func GenProto(projectName string) {

}

func main() {
	cmd.Execute()
	/*	projectName := "surl"
		CreateProject(projectName)*/
}
