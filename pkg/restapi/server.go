package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"strings"
	"tx55/pkg/restapi/events"
	"tx55/pkg/restapi/gameweb"
)

var l = logrus.StandardLogger()

type Server struct {
	DB           *gorm.DB
	Engine       *gin.Engine
	EventService *events.Service
}

func NewServer(db *gorm.DB) *Server {
	gin.SetMode(gin.ReleaseMode)

	s := &Server{
		DB:     db,
		Engine: gin.Default(),
	}

	s.Engine.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("db", s.DB)
			c.Next()
		}
	}())

	if tokens, found := os.LookupEnv("EVENT_TOKENS"); found {
		tokenList := strings.Split(tokens, ",")
		if len(tokenList) > 0 {
			s.EventService = events.NewService(logrus.StandardLogger(), tokenList)
			go s.EventService.Run()
		}
	}

	_ = s.Engine.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	s.Engine.GET("/api/v1/news/list", getNewsList)
	s.Engine.GET("/api/v1/lobby/list", getLobbyList)
	s.Engine.GET("/api/v1/rankings/:period", getRankings)
	s.Engine.GET("/api/v1/rankings/:period/:page", getRankings)
	s.Engine.GET("/api/v1/user/:user_id", getUser)
	s.Engine.GET("/api/v1/user/:user_id/stats", getUserStats)
	s.Engine.GET("/api/v1/user/:user_id/games", getUserGames)
	s.Engine.GET("/api/v1/user/:user_id/games/:page", getUserGames)
	s.Engine.GET("/api/v1/user/:user_id/settings", getUserOptions)
	s.Engine.GET("/api/v1/game/list", getGamesList)
	s.Engine.GET("/api/v1/game/:game_id", getGame)

	if s.EventService != nil {
		s.Engine.POST("/api/v1/stream/events/:token", s.EventService.PostNewEvent)
		s.Engine.GET("/api/v1/stream/events", s.EventService.AcceptGinWebsocket)
	} else {
		s.Engine.GET("/api/v1/stream/events", notImplemented)
		s.Engine.POST("/api/v1/stream/events/:token", notImplemented)
	}

	s.Engine.POST("/us/mgs3/rank/mg3getrank.html", gameweb.PostGetRanks)
	s.Engine.GET("/us/mgs3/text/:filename", gameweb.GetTextFile)
	s.Engine.POST("/us/mgs3/reguser/reguser.html", gameweb.RegisterAccount)
	s.Engine.POST("/us/mgs3/reguser/deluser.html", gameweb.DeleteAccount)
	s.Engine.POST("/us/mgs3/reguser/chgpswd.html", gameweb.ChangePassword)

	return s
}

func success(c *gin.Context, data any) {
	c.JSON(200, ResponseJSON{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, ResponseJSON{
		Success: false,
		Data:    message,
	})
}

func notImplemented(c *gin.Context) {
	success(c, "Not implemented")
}
