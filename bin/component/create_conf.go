package component

import (
	"path/filepath"
	"time"
)

type GatewayConfigParams struct {
	//field for gen
	CreateTime         time.Time
	AppName            string   `json:"app_name"`
	Domain             string   `json:"domain"`
	Version            string   `json:"version"`
	GatewayAddr        string   `json:"gateway_addr"`
	EtcdEnable         bool     `json:"etcd_enable"`
	EtcdServers        []string `json:"etcd_servers"`
	JaegerEnable       bool     `json:"jaeger_enable"`
	JaegerServers      []string `json:"jaeger_servers"`
	KeystorePublicKey  string   `json:"keystore_public_key"`
	KeystorePrivateKey string   `json:"keystore_private_key"`
}

func GatewayConfig(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "gateway", "conf", "app.yaml")
	domain, appName := GetDomainAppName(projectBase)
	c := &GatewayConfigParams{
		CreateTime:         time.Now(),
		Domain:             domain,
		AppName:            appName,
		Version:            "v1.0.0",
		GatewayAddr:        ":8081",
		EtcdEnable:         true,
		EtcdServers:        []string{"http://127.0.0.1:2379"},
		JaegerEnable:       true,
		JaegerServers:      []string{"localhost:6831"},
		KeystorePublicKey:  "../keystore/grpc.pem",
		KeystorePrivateKey: "../keystore/grpc.key",
	}

	tmplStr, err := GetTmpl("gateway_conf_app.yaml.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

type ServerConfigParams struct {
	//field for gen
	CreateTime         time.Time
	Domain             string   `json:"domain"`
	AppName            string   `json:"app_name"`
	Version            string   `json:"version"`
	ServerAddr         string   `json:"server_addr"`
	EtcdEnable         bool     `json:"etcd_enable"`
	EtcdServers        []string `json:"etcd_servers"`
	JaegerEnable       bool     `json:"jaeger_enable"`
	JaegerServers      []string `json:"jaeger_servers"`
	KeystorePublicKey  string   `json:"keystore_public_key"`
	KeystorePrivateKey string   `json:"keystore_private_key"`
}

func ServerConfig(projectBase string) (err error) {
	absFile := filepath.Join(projectBase, "server", "conf", "app.yaml")
	domain, appName := GetDomainAppName(projectBase)
	c := &ServerConfigParams{
		CreateTime:         time.Now(),
		Domain:             domain,
		AppName:            appName,
		Version:            "v1.0.0",
		ServerAddr:         "127.0.0.1:9091",
		EtcdEnable:         true,
		EtcdServers:        []string{"http://127.0.0.1:2379"},
		JaegerEnable:       true,
		JaegerServers:      []string{"localhost:6831"},
		KeystorePublicKey:  "../keystore/grpc.pem",
		KeystorePrivateKey: "../keystore/grpc.key",
	}
	tmplStr, err := GetTmpl("server_conf_app.yaml.tmpl")
	if err != nil {
		return
	}
	err = DoWriteFile(tmplStr, c, absFile, NewDoWriteFileOption())
	return
}

func CreateConf(projectBase string) (err error) {
	GatewayConfig(projectBase)
	ServerConfig(projectBase)
	return
}
