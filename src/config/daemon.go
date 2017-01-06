package config

type Daemon struct {
	dbTransport *DBTransport
	mQ *MQ
}


type DBTransport struct {

}


type MQ struct {

}


func (d *Daemon) initDB(){

}

func (d *Daemon) initMQ(){

}

func InitDaemon() *Daemon{

	d := new(Daemon)
	d.initDB()
	d.initMQ()
	return d
}
