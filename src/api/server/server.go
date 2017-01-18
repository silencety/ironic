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
	"strings"
	"src/github.com/gorilla/mux"
	"api/server/httputils"
	"src/golang.org/x/net/context"
	"utils/errcode"
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
	daemon        *config.Daemon
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
	//s.addRouter(xxx.NewRouter(d))
}

// addRouter adds a new router to the server.
func (s *Server) addRouter(r router.Router) {
	s.routers = append(s.routers, r)
}

// SetDaemon initializes the daemon field
func (s *Server) SetDaemon(daemon *config.Daemon) {
	s.daemon = daemon
}

// Wait blocks the server goroutine until it exits.
// It sends an error message if there is any error during
// the API execution.
func (s *Server) Wait(waitChan chan error) {
	if err := s.serveAPI(); err != nil {
		logrus.Errorf("ServeAPI error: %v", err)
		waitChan <- err
		return
	}
	waitChan <- nil
}

func (s *Server) initRouterSwapper() {
	s.routerSwapper = &routerSwapper{
		router: s.createMux(),
	}
}

// createMux initializes the main router the server uses.
// we keep enableCors just for legacy usage, need to be removed in the future
func (s *Server) createMux() *mux.Router {
	m := mux.NewRouter()
	logrus.Debugf("Registering routers")
	for _, apiRouter := range s.routers {
		for _, r := range apiRouter.Routes() {
			f := s.makeHTTPHandler(r.Handler())

			logrus.Debugf("Registering %s, %s", r.Method(), r.Path())
			//m.Path(versionMatcher + r.Path()).Methods(r.Method()).Handler(f)
			m.Path(r.Path()).Methods(r.Method()).Handler(f)
		}
	}

	return m
}

func (s *Server) makeHTTPHandler(handler httputils.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log the handler call
		logrus.Debugf("Calling %s %s", r.Method, r.URL.Path)

		// Define the context that we'll pass around to share info
		// like the docker-request-id.
		//
		// The 'context' will be used for global data that should
		// apply to all requests. Data that is specific to the
		// immediate function being called should still be passed
		// as 'args' on the function call.
		ctx := context.Background()
		//handlerFunc := s.handleWithGlobalMiddlewares(handler)
		handlerFunc := handler
		vars := mux.Vars(r)
		if vars == nil {
			vars = make(map[string]string)
		}

		if err := handlerFunc(ctx, w, r, vars); err != nil {
			logrus.Errorf("Handler for %s %s returned error: %s", r.Method, r.URL.Path, errcode.GetErrorMessage(err))
			httputils.WriteError(w, err)
		}
	}
}

// serveAPI loops through all initialized servers and spawns goroutine
// with Server method for each. It sets createMux() as Handler also.
func (s *Server) serveAPI() error {
	s.initRouterSwapper()

	var chErrors = make(chan error, len(s.servers))
	for _, srv := range s.servers {
		srv.srv.Handler = s.routerSwapper
		go func(srv *HTTPServer) {
			var err error
			logrus.Infof("API listen on %s", srv.l.Addr())
			if err = srv.Serve(); err != nil && strings.Contains(err.Error(), "use of closed network connection") {
				err = nil
			}
			chErrors <- err
		}(srv)
	}

	for i := 0; i < len(s.servers); i++ {
		err := <-chErrors
		if err != nil {
			return err
		}
	}

	return nil
}

// Serve starts listening for inbound requests.
func (s *HTTPServer) Serve() error {
	return s.srv.Serve(s.l)
}

// Close closes servers and thus stop receiving requests
func (s *Server) Close() {
	for _, srv := range s.servers {
		if err := srv.Close(); err != nil {
			logrus.Error(err)
		}
	}
}

// Close closes the HTTPServer from listening for the inbound requests.
func (s *HTTPServer) Close() error {
	return s.l.Close()
}




