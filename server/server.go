package server

import (
	"crypto/tls"
	"net/http"
	"time"

	logger "log"

	log "github.com/Sirupsen/logrus"
	"github.com/mephux/kolide/config"
)

// Server configurations
type Server struct {
	Addr    string
	Cert    string
	Key     string
	Http    *http.Server
	Version string
}

// Load returns a new server struct based of the given configuration
// struct
func Load(config *config.Config) *Server {
	return &Server{
		Addr: config.Server.Address,
		Cert: config.Server.Crt,
		Key:  config.Server.Key,
	}
}

// Run the server
func (s *Server) Run(handler http.Handler) {

	w := log.StandardLogger().Writer()
	defer w.Close()

	httpServer := &http.Server{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Addr:           s.Addr,
		Handler:        handler,
		ReadTimeout:    2 * time.Minute,
		WriteTimeout:   2 * time.Minute,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       logger.New(w, "", 0),
		// ConnState: func(connection net.Conn, state http.ConnState) {
		// log.Debugf("Connection State Change: %s", state.String())
		// },
	}

	if len(s.Cert) != 0 {
		log.Infof("Starting server: https://%s", s.Addr)
		log.Fatal(httpServer.ListenAndServeTLS(s.Cert, s.Key))
	} else {
		log.Infof("Starting server: http://%s", s.Addr)
		log.Fatal(httpServer.ListenAndServe())
	}
}
