package api

import (
	"fmt"
	"io"
	"tpi-mon/tpi"

	"github.com/gin-gonic/gin"
)

// Run the api with supplied tpi, and binding to supplied port
func Run(tpiClient tpi.Client, bindPort uint16) error {
	r := gin.Default()
	setupRoutes(r, tpiClient)
	bindAddr := fmt.Sprintf("0.0.0.0:%d", bindPort)
	return r.Run(bindAddr)
}

func setupRoutes(r *gin.Engine, tpiClient tpi.Client) {

	r.GET("/", func(c *gin.Context) {
		c.String(200, "tpimon api 1.0")
	})

	r.GET("/partitions", func(c *gin.Context) {
		c.JSON(200, tpiClient.GetPartitions())
	})

	r.GET("/zones", func(c *gin.Context) {
		c.JSON(200, tpiClient.GetZones())
	})

	r.GET("/events", func(c *gin.Context) {
		// note this endpoint should not actually be the events ch consumer
		eventsCh := tpiClient.GetEventCh()
		c.Stream(func(w io.Writer) bool {
			c.SSEvent("event", <-eventsCh)
			return true
		})
	})

	r.POST("/partitions/:partID/arm/:type", func(c *gin.Context) {

		partID := c.Param("partID")
		armType := c.Param("type")

		switch armType {
		case "away":
			tpiClient.AwayArm(partID)
		case "stay":
			tpiClient.StayArm(partID)
		case "zero-delay":
			tpiClient.ZeroEntryDelayArm(partID)
		case "with-code":
			data := struct {
				PIN string `json:"pin" binding:"required"`
			}{}

			if err := c.BindJSON(&data); err != nil {
				c.JSON(400, &gin.H{"error": err.Error()})
				return
			}

			tpiClient.ArmWithCode(partID, data.PIN)
		}

		c.JSON(202, nil)
	})

	r.POST("/partitions/:partID/disarm", func(c *gin.Context) {

		partID := c.Param("partID")
		data := struct {
			PIN string `json:"pin" binding:"required"`
		}{}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
			return
		}

		tpiClient.Disarm(partID, data.PIN)
		c.JSON(202, nil)
	})

}
