package config

import (
	"fmt"
	"os"
	"src/github.com/ini"

)



type Global struct {
	Platform 	string
	Name 		string
}




var (
	GlobalConfig *Global
	DataBaseConfig *Database
	configFileName string
)


func ConfigInit(file string) (cfg *ini.File){


	GlobalConfig = &Global{Platform:"linux", Name:"global"}
	cfg, err := ini.Load(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//
	//bind_addr := cfg.Section("default").Key("bind").String()
	//fmt.Println(bind_addr)
	//DataBaseConfig = new(Database)
	//fmt.Println(cfg.Section("database").Key("connection").String())
	//err = cfg.Section("database").MapTo(DataBaseConfig)
	//if err == nil {
	//	fmt.Println(DataBaseConfig.Connection)
	//}else {
	//	fmt.Println(err)
	//}
        return  cfg

}
