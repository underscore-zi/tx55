package restapi

import (
	"errors"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"tx55/pkg/configurations"
	"tx55/pkg/restapi/events"
	"tx55/pkg/restapi/gameweb"
)

type Server struct {
	DB           *gorm.DB
	Engine       *gin.Engine
	EventService *events.Service
}

func NewServer(config configurations.RestAPI) (s *Server, err error) {
	gin.SetMode(gin.ReleaseMode)
	s = &Server{
		Engine: gin.Default(),
	}

	s.DB, err = config.Database.Open(&gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", 0), config.Database.LogConfig.LoggerConfig()),
	})
	if err != nil {
		return
	}

	if config.SessionSecret == "" {
		err = errors.New("no session secret specified")
		return
	}

	sessionDB, err := config.SessionDatabase.Open(&gorm.Config{})
	if err != nil {
		return
	}
	store := gormsessions.NewStore(sessionDB, true, []byte(config.SessionSecret))
	s.Engine.Use(sessions.Sessions("sessions", store))

	s.Engine.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("logger", logrus.StandardLogger())
			c.Next()
		}
	}())

	s.Engine.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("db", s.DB)
			c.Next()
		}
	}())

	if config.Events.Enabled {
		if len(config.Events.AccessTokens) == 0 {
			err = errors.New("no event access tokens specified")
			return
		}
		s.EventService = events.NewService(logrus.StandardLogger(), config.Events.AccessTokens)
		go s.EventService.Run()
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

	s.Engine.POST("/api/v1/auth/login", Login)
	s.Engine.GET("/api/v1/auth/logout", Logout)

	gameusers := s.Engine.Group("/api/v1")
	gameusers.Use(GameLoginRequired)
	gameusers.GET("/user/whoami", whoAmI)
	gameusers.POST("/user/display_name")
	return
}

func success(c *gin.Context, data any) {
	c.JSON(200, ResponseJSON{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, ResponseJSON{
		Success: false,
		Data:    message,
	})
}

func notImplemented(c *gin.Context) {
	success(c, "Not implemented")
}
