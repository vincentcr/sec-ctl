package main

import (
	"fmt"
	"io"
	"sec-ctl/pkg/sites"

	"github.com/gin-gonic/gin"
)

//type authenticate func(clientID string, secret string) (sites.Site, bool)

// run starts the api with supplied tpi, and binding to supplied prt
func runRESTAPI(site sites.Site, bindHost string, bindPort uint16) error {
	g := gin.Default()
	setupRoutes(site, g)
	bindAddr := fmt.Sprintf("%s:%d", bindHost, bindPort)
	return g.Run(bindAddr)
}

func setupRoutes(site sites.Site, g *gin.Engine) {

	g.GET("/", func(c *gin.Context) {
		c.JSON(200, site.GetState())
	})

	g.POST("/commands", func(c *gin.Context) {

		var cmd sites.UserCommand
		if err := c.BindJSON(&cmd); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
			return
		}

		if err := site.Exec(cmd); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
			return
		}

		c.JSON(202, "Command sent")
	})

	g.GET("/events", func(c *gin.Context) {

		eventCh := site.SubscribeToEvents()
		c.Stream(func(w io.Writer) bool {
			c.SSEvent("event", <-eventCh)
			return true
		})
	})

}
