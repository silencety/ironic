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
	"runtime/debug"
)

const defaultDaemonConfigFile = "/root/ironic/conf/ironic.conf"
const defaultLogFile = "/var/log/ironic.log"

var(
	configFileName string
	GlobalConfig *ini.File
)


func init() {
	flag.StringVar(&configFileName, "config", defaultDaemonConfigFile, "configure file!")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION] \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	defer func(){
		if err :=recover(); err != nil {
			debug.PrintStack()
		}
	}()
        flag.Parse()
	GlobalConfig = config.ConfigInit(getDefaultConfigFile())
	initLog()
	initServer()

}


func getDefaultConfigFile() string {
	return configFileName
}


// init logrus: set formatter
//              set out put
//              set log level
//		create a goroutine to set log level by signal:
func initLog(){

	logFile := GlobalConfig.Section(config.SectionDefault).Key("logFile").String()
	if logFile == ""{
		logFile = defaultLogFile
	}
	f, err := os.OpenFile(logFile, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logrus.SetOutput(f)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: jsonlog.RFC3339NanoFixed,
		FullTimestamp: true,
	})
	setDaemonLogLevel(GlobalConfig.Section(config.SectionDefault).Key("loglevel").String())
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


