package main

import (
	"api/server"
	"api"
	"src/github.com/Sirupsen/logrus"
	"fmt"
	"config"
)

const defaultBindAddr = "0.0.0.0"
const defaultbindPort = 8778

//read config from global config
//init http server
var d *config.Daemon


func GetServer(){
	serverConfig := &server.Config{
		Logging:	true,
		Version:	api.Version(),
	}
	bind_addr := GlobalConfig.Section(sectionDefault).Key("bind").String()
	if bind_addr == "" {
		bind_addr = defaultBindAddr
	}

	bind_port,err := GlobalConfig.Section(sectionDefault).Key("port").Int()
	if err != nil {
		bind_port = defaultbindPort
	}
	serverConfig.Addrs = append(serverConfig.Addrs,
		server.Addr{ Proto: "tcp", Addr: fmt.Sprintf("%s:%v", bind_addr, bind_port)})
	logrus.Debugf("read socket form config file  %s", serverConfig.Addrs[0].Addr)
	api, err := server.New(serverConfig)
	if err != nil {
		logrus.Fatal(err)
	}
	d = config.InitDaemon()
	api.InitRouters(d)

}
