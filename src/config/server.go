package config


type Server struct {
	Logging                  bool
	EnableCors               bool
	CorsHeaders              string
	Version                  string
	Addrs                    []Addr
}

var ServerCfg = &Server{}


func init(){
	addr := Addr{Proto : "tcp", Addr:"127.0.0.1:9090"}
	_ = append(ServerCfg.Addrs, addr)
}



