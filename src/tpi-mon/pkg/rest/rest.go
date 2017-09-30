package rest

import (
	"fmt"
	"io"
	"tpi-mon/pkg/site"

	"github.com/gin-gonic/gin"
)

type lookupClient func(id string) site.Client

// Start stats the api with supplied tpi, and binding to supplied port
func Start(lookupClient lookupClient, bindHost string, bindPort uint16, errCh chan error) {
	g := gin.Default()
	setupRoutes(g, lookupClient)
	bindAddr := fmt.Sprintf("%s:%d", bindHost, bindPort)
	go func() {
		err := g.Run(bindAddr)
		if err != nil {
			errCh <- err
		}
	}()
}

func setupRoutes(g *gin.Engine, lookupClient lookupClient) {

	g.GET("/", func(c *gin.Context) {
		c.String(200, "tpimon api 1.0")
	})

	g.GET("/clients/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, lookupClient(id).GetState())
	})

	g.POST("/clients/:id/commands", func(c *gin.Context) {
		id := c.Param("id")

		var cmd site.UserCommand
		if err := c.BindJSON(&cmd); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
			return
		}

		if err := lookupClient(id).Exec(cmd); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
			return
		}

		c.JSON(202, "")
	})

	g.GET("/clients/:id/events", func(c *gin.Context) {
		id := c.Param("id")
		eventCh := lookupClient(id).SubscribeToEvents()
		c.Stream(func(w io.Writer) bool {
			c.SSEvent("event", <-eventCh)
			return true
		})
	})
}
