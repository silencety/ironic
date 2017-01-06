package config

type Database struct {

	Connection string 	`ini:"connection"`
	MaxTimeOut string	`ini:"max_time_out"`
}
