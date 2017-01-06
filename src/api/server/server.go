package server

import (
	"net/http"
	"net"
	"api/server/router"
	"src/github.com/Sirupsen/logrus"
	"fmt"
	"crypto/tls"
	"utils/socket"
	"utils/portallocator"
	"strconv"
	"api/server/router/host"
	"config"
)

type Config struct {
	Logging                  bool
	EnableCors               bool
	CorsHeaders              string
	//AuthorizationPluginNames []string
	Version                  string
	//SocketGroup              string
	TLSConfig                *tls.Config
	Addrs                    []Addr
}

type Addr struct {
	Proto string
	Addr  string
}

type Server struct {
	cfg           *Config
	servers       []*HTTPServer
	routers       []router.Router
	//authZPlugins  []authorization.Plugin
	routerSwapper *routerSwapper
	//daemon        *daemon.Daemon
}

type HTTPServer struct {
	srv *http.Server
	l   net.Listener
}

func New(cfg *Config) (*Server, error) {
	s := &Server{
		cfg: cfg,
	}
	for _, addr := range cfg.Addrs {
		srv, err := s.newServer(addr.Proto, addr.Addr)
		if err != nil {
			return nil, err
		}
		logrus.Debugf("Server created for HTTP on %s (%s)", addr.Proto, addr.Addr)
		s.servers = append(s.servers, srv...)
	}
	return s, nil
}

func (s *Server) newServer(proto, addr string) ([]*HTTPServer, error) {
	var (
		ls  []net.Listener
	)
	switch proto {
	case "tcp":
		l, err := s.initTCPSocket(addr)
		if err != nil {
			return nil, err
		}
		ls = append(ls, l)
	ls = append(ls, l)
	default:
		return nil, fmt.Errorf("Invalid protocol format: %q", proto)
	}
	var res []*HTTPServer
	for _, l := range ls {
		res = append(res, &HTTPServer{
			&http.Server{
				Addr: addr,
			},
			l,
		})
	}
	return res, nil
}


func (s *Server) initTCPSocket(addr string) (l net.Listener, err error) {
	if s.cfg.TLSConfig == nil || s.cfg.TLSConfig.ClientAuth != tls.RequireAndVerifyClientCert {
		logrus.Warn("/!\\ DON'T BIND ON ANY IP ADDRESS WITHOUT setting -tlsverify IF YOU DON'T KNOW WHAT YOU'RE DOING /!\\")
	}
	if l, err = socket.NewTCPSocket(addr, s.cfg.TLSConfig); err != nil {
		return nil, err
	}
	if err := allocateDaemonPort(addr); err != nil {
		return nil, err
	}
	return
}

func allocateDaemonPort(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	var hostIPs []net.IP
	if parsedIP := net.ParseIP(host); parsedIP != nil {
		hostIPs = append(hostIPs, parsedIP)
	} else if hostIPs, err = net.LookupIP(host); err != nil {
		return fmt.Errorf("failed to lookup %s address in host specification", host)
	}

	pa := portallocator.Get()
	for _, hostIP := range hostIPs {
		if _, err := pa.RequestPort(hostIP, "tcp", intPort); err != nil {
			return fmt.Errorf("failed to allocate daemon listening port %d (err: %v)", intPort, err)
		}
	}
	return nil
}

// InitRouters initializes a list of routers for the server.
func (s *Server) InitRouters(d *config.Daemon) {
	s.addRouter(host.NewRouter(d))
}

// addRouter adds a new router to the server.
func (s *Server) addRouter(r router.Router) {
	s.routers = append(s.routers, r)
}





