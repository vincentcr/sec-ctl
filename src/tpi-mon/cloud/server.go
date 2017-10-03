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
}

func startServer(bindHost string, bindPort uint16) *server {
	s := &server{
		registry: newRegistry(),
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
		if err := s.wshandler(c.Writer, c.Request); err != nil {
			logger.Println("Unable to upgrade request to websocket:", err)
			c.JSON(400, &gin.H{"error": "Unable to upgrade to web socket"})
		}
	})

	err := r.Run(addr)
	if err != nil {
		logger.Panicln(err)
	}
}

func (s *server) wshandler(w http.ResponseWriter, r *http.Request) error {

	conn, err := ws.UpgradeRequest(w, r)
	if err != nil {
		return err
	}

	initRemoteSite(conn, s.registry)
	return nil
}

func (s *server) GetClient(id string) site.Client {
	return s.registry.getClient(id)
}
