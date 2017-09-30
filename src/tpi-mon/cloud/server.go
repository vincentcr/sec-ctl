package main

import (
	"fmt"
	"net/http"
	"tpi-mon/pkg/site"
	"tpi-mon/pkg/ws"

	"github.com/gin-gonic/gin"
)

type server struct {
	registry *clientRegistry
	errCh    chan error
}

func startServer(bindHost string, bindPort uint16, errCh chan error) *server {
	s := &server{
		registry: newRegistry(),
		errCh:    errCh,
	}
	go func() {
		s.runWSServer(bindHost, bindPort)
	}()
	return s
}

func (s *server) runWSServer(bindHost string, bindPort uint16) {
	r := gin.Default()

	addr := fmt.Sprintf("%s:%d", bindHost, bindPort)

	r.GET("/ws", func(c *gin.Context) {
		s.wshandler(c.Writer, c.Request)
	})

	err := r.Run(addr)
	if err != nil {
		s.errCh <- err
	}
}

func (s *server) wshandler(w http.ResponseWriter, r *http.Request) {

	conn, err := ws.UpgradeRequest(w, r)
	if err != nil {
		s.errCh <- err
	}

	initRemoteSite(conn, s.registry)
}

func (s *server) GetClient(id string) site.Client {
	return s.registry.getClient(id)
}
