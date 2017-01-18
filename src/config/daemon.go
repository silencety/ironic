package config

import "src/github.com/ini"

type Daemon struct {
	dbTransport *DBTransport
	mQ *MQ
	PidFile string
}


type DBTransport struct {

}


type MQ struct {

}


func (d *Daemon) initDB(GlobalConfig *ini.File){

}


func (d *Daemon) initMQ(GlobalConfig *ini.File){

}

func InitDaemon(GlobalConfig *ini.File) *Daemon{

	d := new(Daemon)
	d.initDB(GlobalConfig)
	d.initMQ(GlobalConfig)
	pidFile := GlobalConfig.Section(SectionDefault).Key("pidFile").String()
	if pidFile == "" {
		pidFile = defaultPidFile
	}
	d.PidFile = pidFile
	return d
}

func (d *Daemon) Shutdown(){
}