package main

import (
	"database/sql"
	"fmt"
	"strings"
	"tpi-mon/cloud/db"
	"tpi-mon/pkg/sites"
	"tpi-mon/pkg/ws"

	"github.com/gin-gonic/gin"
)

type rest struct {
	db       *db.DB
	registry *siteRegistry
	gin      *gin.Engine
}

func runRESTAPI(reg *siteRegistry, db *db.DB, bindHost string, bindPort uint16) error {
	rest := rest{
		gin:      gin.Default(),
		registry: reg,
		db:       db,
	}

	rest.setup()
	bindAddr := fmt.Sprintf("%s:%d", bindHost, bindPort)
	return rest.gin.Run(bindAddr)
}

func (rest rest) setup() {

	rest.gin.GET("/", func(c *gin.Context) {
		c.String(200, "tpimon api 1.0")
	})

	rest.gin.GET("/ws", rest.authSiteByToken(), func(c *gin.Context) {
		conn, err := ws.UpgradeRequest(c.Writer, c.Request)
		if err != nil {
			logger.Println("Unable to upgrade request to websocket:", err)
			c.JSON(400, &gin.H{"error": "Unable to upgrade to web socket"})
		}

		site := c.MustGet("Site").(db.Site)
		rest.registry.initRemoteSite(site, conn)
	})

	rest.gin.POST("/signup", func(c *gin.Context) {
		var userForm struct {
			Email    string `binding:"required"`
			Password string `binding:"required"`
		}

		if err := c.BindJSON(&userForm); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
			return
		}
		user, tok, err := rest.db.CreateUser(userForm.Email, userForm.Password)
		if err != nil {
			c.JSON(500, &gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, &gin.H{
			"user":  user,
			"token": tok,
		})
	})

	rest.gin.POST("/sites", func(c *gin.Context) {

		site, tok, claimTok, err := rest.db.CreateSite()
		if err != nil {
			c.JSON(500, &gin.H{"error": err.Error()})
			return
		}

		setupURL := fmt.Sprintf("/sites/%s/firstTime?t=%s", site.ID, claimTok)

		c.JSON(200, &gin.H{
			"SiteID":   site.ID,
			"Token":    tok,
			"SetupURL": setupURL,
		})
	})

	rest.gin.POST("/sites/:id/claim", rest.authUserByToken(), func(c *gin.Context) {
		site := c.MustGet("Site").(*remoteSite)
		user := c.MustGet("User").(db.User)

		var claimForm struct {
			ClaimToken string `binding:"required"`
		}
		if err := c.BindJSON(&claimForm); err != nil {
			c.JSON(400, &gin.H{"error": err.Error()})
			return
		}

		rest.db.ClaimSite(user, site.id, claimForm.ClaimToken)
	})

	sitesRouter := rest.gin.Group("/sites/:id", rest.authUserByToken(), rest.getSite())
	{
		sitesRouter.GET("/", func(c *gin.Context) {
			site := c.MustGet("Site").(*remoteSite)
			c.JSON(200, site.GetState())
		})

		sitesRouter.POST("/commands", func(c *gin.Context) {

			var cmd sites.UserCommand
			if err := c.BindJSON(&cmd); err != nil {
				c.JSON(400, &gin.H{"error": err.Error()})
				return
			}

			site := c.MustGet("Site").(db.Site)
			if err := rest.registry.sendCommand(site.ID, cmd); err != nil {
				c.JSON(400, &gin.H{"error": err.Error()})
				return
			}

			c.JSON(202, "Command sent")
		})

		sitesRouter.GET("/events", func(c *gin.Context) {

			site := c.MustGet("Site").(db.Site)
			evts, err := rest.registry.getLatestEvents(site.ID, 100)
			if err != nil {
				c.JSON(500, "Internal Error")
				return
			}

			c.JSON(200, evts)

		})
	}
}

// func setupRoutes(g *gin.Engine, reg *siteRegistry, db *db.DB) {

// 	g.GET("/", func(c *gin.Context) {
// 		c.String(200, "tpimon api 1.0")
// 	})

// 	g.GET("/ws", authSiteByToken(db), func(c *gin.Context) {
// 		conn, err := ws.UpgradeRequest(c.Writer, c.Request)
// 		if err != nil {
// 			logger.Println("Unable to upgrade request to websocket:", err)
// 			c.JSON(400, &gin.H{"error": "Unable to upgrade to web socket"})
// 		}

// 		site := c.MustGet("Site").(db.Site)
// 		initRemoteSite(conn, reg, site)
// 	})

// 	g.POST("/signup", func(c *gin.Context) {
// 		var userForm struct {
// 			Email    string `binding:"required"`
// 			Password string `binding:"required"`
// 		}

// 		if err := c.BindJSON(&userForm); err != nil {
// 			c.JSON(400, &gin.H{"error": err.Error()})
// 			return
// 		}
// 		user, tok, err := db.CreateUser(userForm.Email, userForm.Password)
// 		if err != nil {
// 			c.JSON(500, &gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(200, &gin.H{
// 			"user":  user,
// 			"token": tok,
// 		})
// 	})

// 	g.POST("/sites", func(c *gin.Context) {

// 		site, tok, claimTok, err := db.CreateSite()
// 		if err != nil {
// 			c.JSON(500, &gin.H{"error": err.Error()})
// 			return
// 		}

// 		setupURL := fmt.Sprintf("/sites/%s/firstTime?t=%s", site.ID, claimTok)

// 		c.JSON(200, &gin.H{
// 			"SiteID":   site.ID,
// 			"Token":    tok,
// 			"SetupURL": setupURL,
// 		})
// 	})

// 	g.POST("/sites/:id/claim", authUserByToken(db), func(c *gin.Context) {
// 		site := c.MustGet("Site").(*remoteSite)
// 		user := c.MustGet("User").(db.User)

// 		var claimForm struct {
// 			ClaimToken string `binding:"required"`
// 		}
// 		if err := c.BindJSON(&claimForm); err != nil {
// 			c.JSON(400, &gin.H{"error": err.Error()})
// 			return
// 		}

// 		db.ClaimSite(user, site.id, claimForm.ClaimToken)
// 	})

// 	sitesRouter := g.Group("/sites/:id", authUserByToken(db), getSite(reg, db))
// 	{
// 		sitesRouter.GET("/", func(c *gin.Context) {
// 			site := c.MustGet("Site").(*remoteSite)
// 			c.JSON(200, site.GetState())
// 		})

// 		sitesRouter.POST("/commands", func(c *gin.Context) {
// 			site := c.MustGet("Site").(*remoteSite)

// 			var cmd sites.UserCommand
// 			if err := c.BindJSON(&cmd); err != nil {
// 				c.JSON(400, &gin.H{"error": err.Error()})
// 				return
// 			}

// 			if err := site.Exec(cmd); err != nil {
// 				c.JSON(400, &gin.H{"error": err.Error()})
// 				return
// 			}

// 			c.JSON(202, "Command sent")
// 		})

// 		sitesRouter.GET("/events", func(c *gin.Context) {
// 			site := c.MustGet("Site").(*remoteSite)

// 			eventCh := site.SubscribeToEvents()
// 			c.Stream(func(w io.Writer) bool {
// 				c.SSEvent("event", <-eventCh)
// 				return true
// 			})
// 		})
// 	}
// }

func requestAuthToken(c *gin.Context) string {

	// try _t query param
	tok := c.Query("_t")
	if tok != "" {
		return tok
	}

	// try Authorization: Bearer <tok> header
	hdPfx := "Bearer "
	hdVal := c.GetHeader("authorization")
	if strings.HasPrefix(hdVal, hdPfx) {
		return hdVal[len(hdPfx):]
	}

	// try token cookie
	tok, _ = c.GetCookie("token")
	if tok != "" {
		return tok
	}

	return ""
}

func authResourceByToken(key string, fetch func(string) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {

		tok := requestAuthToken(c)
		if tok == "" {
			c.AbortWithStatusJSON(401, &gin.H{"error": "Authorization required"})
		} else {
			res, err := fetch(tok)
			if err == nil {
				c.Set(key, res)
			} else if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(401, &gin.H{"error": "Invalid authorization token"})
			} else {
				logger.Println("unexpected error", err)
				c.AbortWithStatus(500)
			}
		}
	}
}

func (rest rest) authUserByToken() gin.HandlerFunc {
	return authResourceByToken("User", func(token string) (interface{}, error) {
		return rest.db.AuthUserByToken(token)
	})
}

func (rest rest) authSiteByToken() gin.HandlerFunc {
	return authResourceByToken("Site", func(token string) (interface{}, error) {
		return rest.db.AuthSiteByToken(token)
	})
}

func (rest rest) getSite() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := c.MustGet("User").(db.User)

		id := db.UUID(c.Param("id"))
		site, err := rest.registry.getSite(user, id)

		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatus(404)
			} else if err.Error() == "Unauthorized" {
				c.AbortWithStatus(403)
			} else {
				logger.Printf("Error fetching site: %v\n", err)
				c.AbortWithStatus(500)
			}
			return
		}

		c.Set("Site", site)

	}
}
