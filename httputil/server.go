package httputil

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mdzio/go-logging"
)

const shutdownTimeout = 5 * time.Second

// Server is a plain http.Server with some usefull additions. The
// http.DefaultServeMux is used for the handlers.
type Server struct {
	// Binding address for serving HTTP.
	Addr string
	// Binding address for serving HTTPS.
	AddrTLS string
	// Certificate file for HTTPS
	CertFile string
	// Private key file for HTTPS
	KeyFile string
	// When an error happens while serving (e.g. binding of port fails), this
	// error is sent to the channel ServeErr.
	ServeErr chan<- error
	// Default logger is "http-server", if not specified
	Log logging.Logger

	done      chan struct{}
	server    http.Server
	serverTLS http.Server
}

// Startup starts the HTTP server.
func (s *Server) Startup() {
	// setup
	s.server.Addr = s.Addr
	s.serverTLS.Addr = s.AddrTLS
	// capacity of 2 to avoid blocking, when shutting down
	s.done = make(chan struct{}, 2)
	if s.Log == nil {
		s.Log = logging.Get("http-server")
	}

	// start servers
	if s.server.Addr != "" {
		s.startupServer("HTTP", &s.server, func() error {
			return s.server.ListenAndServe()
		})
	}
	if s.serverTLS.Addr != "" {
		s.startupServer("HTTPS", &s.serverTLS, func() error {
			return s.serverTLS.ListenAndServeTLS(s.CertFile, s.KeyFile)
		})
	}
}

// Shutdown shuts the HTTP server down.
func (s *Server) Shutdown() {
	if s.server.Addr != "" {
		s.shutdownServer("HTTP", &s.server)
	}
	if s.serverTLS.Addr != "" {
		s.shutdownServer("HTTPS", &s.serverTLS)
	}
}

func (s *Server) startupServer(name string, svr *http.Server, runFunc func() error) {
	// start http/s server
	go func() {
		s.Log.Infof("Starting %s server on address %s", name, svr.Addr)
		err := runFunc()
		// signal server is down (must not block)
		s.done <- struct{}{}
		// check for error
		if err != http.ErrServerClosed {
			// signal error while serving (block does not harm)
			if s.ServeErr != nil {
				s.ServeErr <- fmt.Errorf("Running %s server failed: %v", name, err)
			}
		}
	}()
}

func (s *Server) shutdownServer(name string, svr *http.Server) {
	// start shutdown
	s.Log.Debugf("Shutting down %s server", name)
	svr.SetKeepAlivesEnabled(false)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err := svr.Shutdown(ctx)
	if err != nil {
		s.Log.Errorf("Shutdown of %s server failed: %v", name, err)
		return
	}
	// wait for shutdown
	<-s.done
}
