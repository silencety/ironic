package main

import (
	"api/server"
	"api"
	"src/github.com/Sirupsen/logrus"
	"fmt"
	"utils/signal"

	"utils/pidfile"
	"config"
	"time"
)

const defaultBindAddr = "0.0.0.0"
const defaultbindPort = 8778


var d *config.Daemon

//read config from global config
//init http server
func initServer() {
	serverConfig := &server.Config{
		Logging:	true,
		Version:	api.Version(),
	}
	bind_addr := GlobalConfig.Section(config.SectionDefault).Key("bind").String()
	if bind_addr == "" {
		bind_addr = defaultBindAddr
	}

	bind_port, err := GlobalConfig.Section(config.SectionDefault).Key("port").Int()
	if err != nil {
		bind_port = defaultbindPort
	}
	serverConfig.Addrs = append(serverConfig.Addrs,
		server.Addr{ Proto: "tcp", Addr: fmt.Sprintf("%s:%v", bind_addr, bind_port)})
	logrus.Debugf("Read socket form config file  %s", serverConfig.Addrs[0].Addr)
	api, err := server.New(serverConfig)
	if err != nil {
		logrus.Fatal(err)
	}
	var pfile *pidfile.PIDFile
	d = config.InitDaemon(GlobalConfig)
	pf, err := pidfile.New(d.PidFile)
	if err != nil {
		logrus.Fatalf("Error starting daemon: %v", err)
	}
	pfile = pf
	defer func() {
		if err := pfile.Remove(); err != nil {
			logrus.Error(err)
		}
	}()
	api.SetDaemon(d)
	api.InitRouters(d)
	serveAPIWait := make(chan error)
	go api.Wait(serveAPIWait)

	signal.Trap(func() {
		api.Close()
		<-serveAPIWait
		shutdownDaemon(d, 15)
		if err := pfile.Remove(); err != nil {
			logrus.Error(err)
		}
	})

	errAPI := <-serveAPIWait
	shutdownDaemon(d, 15)
	if errAPI != nil {
		if pfile != nil {
			if err := pfile.Remove(); err != nil {
				logrus.Error(err)
			}
		}
		logrus.Fatalf("Shutting down due to ServeAPI error: %v", errAPI)
	}
}


// shutdownDaemon just wraps daemon.Shutdown() to handle a timeout in case
// d.Shutdown() is waiting too long to kill container or worst it's
// blocked there
func shutdownDaemon(d *config.Daemon, timeout time.Duration) {
	ch := make(chan struct{})
	go func() {
		d.Shutdown()
		close(ch)
	}()
	select {
	case <-ch:
		logrus.Debug("Clean shutdown succeeded")
	case <-time.After(timeout * time.Second):
		logrus.Error("Force shutdown daemon")
	}
}
