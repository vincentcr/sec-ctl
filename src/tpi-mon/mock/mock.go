package mock

import (
	"fmt"
	"log"
	"net"
	"os"
	"tpi-mon/tpi"

	"github.com/gin-gonic/gin"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "[mock server] ", log.LstdFlags)
}

// Run creates a new mock server
func Run(bindHost string, tpiPort uint, httpPort uint, stateFile string) error {

	state, err := newState(stateFile)
	if err != nil {
		return err
	}

	ctrl := newController(state)

	errCh := make(chan error)

	startMockTPI(ctrl, bindHost, tpiPort, errCh)
	startRestAPI(ctrl, bindHost, httpPort, errCh)

	for {
		select {
		case err := <-errCh:
			logger.Panicln(err)
		}
	}
}

func startMockTPI(ctrl *controller, bindHost string, bindPort uint, errCh chan error) {
	addr := fmt.Sprintf("%s:%d", bindHost, bindPort)
	logger.Println("binding:", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		errCh <- fmt.Errorf("Error listening: %v", err)
		return
	}

	go func() {
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				errCh <- err
			}
			go handleClientSession(ctrl, conn)
		}
	}()

}

func startRestAPI(ctrl *controller, bindHost string, bindPort uint, errCh chan error) {
	r := gin.Default()
	setupRoutes(r, ctrl)
	go func() {
		if err := r.Run(fmt.Sprintf("%s:%d", bindHost, bindPort)); err != nil {
			errCh <- err
		}
	}()
}

func setupRoutes(r *gin.Engine, ctrl *controller) {
	r.GET("/state", func(c *gin.Context) {
		c.JSON(200, ctrl.state)

	})

	//simulate alarms
	r.POST("/sim/alarms/trigger", func(c *gin.Context) {
		data := struct {
			Type        tpi.AlarmType `json:"type" binding:"required"`
			PartitionID string        `json:"partitionID"`
			ZoneID      string        `json:"zoneID"`
		}{}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
		}

		if err := ctrl.triggerAlarm(data.Type, data.PartitionID, data.ZoneID); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
		} else {
			c.JSON(204, nil)
		}

	})

	//simulate alarm restoral
	r.POST("/sim/alarms/restore", func(c *gin.Context) {
		data := struct {
			Type        tpi.AlarmType `json:"type" binding:"required"`
			PartitionID string        `json:"partitionID"`
		}{}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
			return
		}

		if err := ctrl.restoreAlarm(data.Type, data.PartitionID); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
		}

	})

	//POST sim trouble
	//POST sim trouble restore

}
