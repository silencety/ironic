package main

import (
	"fmt"
	"config"
	"src/github.com/Sirupsen/logrus"
	"src/github.com/jsonlog"
	"os"
	"flag"
	"src/github.com/ini"
	"os/signal"
	"syscall"
)

const defaultDaemonConfigFile = "d:/ironic/conf/ironic.conf"
//const defaultDaemonConfigFile = "/root/ironic/conf/ironic111.conf"

var(
	configFileName string
	GlobalConfig *ini.File
)


func init() {
	flag.StringVar(&configFileName, "config", defaultDaemonConfigFile, "configure file!")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
        flag.Parse()
	GlobalConfig = config.ConfigInit(getDefaultConfigFile())
	initLog()
	GetServer()



}


func getDefaultConfigFile() string {
	return configFileName
}


// init logrus: set formatter
//              set out put
//              set log level
//		support set log level by signal:
func initLog(){
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: jsonlog.RFC3339NanoFixed})
	logrus.SetOutput(os.Stderr)
	setDaemonLogLevel(GlobalConfig.Section(sectionDefault).Key("loglevel").String())
	go func() {
		c := make(chan os.Signal, 10)
		signal.Notify(c, syscall.SIGUSR1,syscall.SIGUSR2)
		for {
			s :=<-c
			switch s {
			case syscall.SIGUSR1:
				logrus.SetLevel(logrus.InfoLevel)
			case syscall.SIGUSR2:
				logrus.SetLevel(logrus.DebugLevel)
			}
		}
	}()
}


