package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var l = logrus.StandardLogger()

type Server struct {
	DB     *gorm.DB
	Engine *gin.Engine
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
	_ = s.Engine.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	s.Engine.GET("/api/v1/lobby/list", getLobbyList)
	s.Engine.GET("/api/v1/user/:user_id", getUser)
	s.Engine.GET("/api/v1/user/:user_id/stats", getUserStats)
	s.Engine.GET("/api/v1/game/list", getGamesList)
	s.Engine.GET("/api/v1/game/:game_id", getGame)

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
