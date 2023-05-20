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
	AdminDB      *gorm.DB
	Engine       *gin.Engine
	EventService *events.Service
}

type AuthLevel int

const (
	AuthLevelNone  AuthLevel = 1
	AuthLevelUser  AuthLevel = 2
	AuthLevelAdmin AuthLevel = 3
)

type endpointInfo struct {
	level            AuthLevel
	method, endpoint string
	handler          gin.HandlerFunc
	args             any
	returns          any
}

var routes = map[AuthLevel][]endpointInfo{
	AuthLevelNone:  {},
	AuthLevelUser:  {},
	AuthLevelAdmin: {},
}

func Register(level AuthLevel, method, endpoint string, handler gin.HandlerFunc) {
	routes[level] = append(routes[level], endpointInfo{
		level:    level,
		method:   method,
		endpoint: endpoint,
		handler:  handler,
	})
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

	s.AdminDB, err = config.AdminDatabase.Open(&gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", 0), config.AdminDatabase.LogConfig.LoggerConfig()),
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
	s.Engine.Use(CORSMiddleware(config))
	s.Engine.Use(sessions.Sessions("sessions", store))
	engineLogger := logrus.StandardLogger()
	s.Engine.Use(ProvideContextVar("logger", engineLogger))
	s.Engine.Use(ProvideContextVar("db", s.DB))
	s.Engine.Use(ProvideContextVar("adminDB", s.AdminDB))

	_ = s.Engine.SetTrustedProxies(config.TrustedProxies)

	unauthGroup := s.Engine.Group(config.ApiPrefix)

	userGroup := s.Engine.Group(config.ApiPrefix, RequireAPIKey, GameLoginRequired)
	adminGroup := s.Engine.Group(config.ApiPrefix, RequireAPIKey, AdminLoginRequired)

	if config.Events.Enabled {
		if len(config.Events.AccessTokens) == 0 {
			err = errors.New("no event access tokens specified")
			return
		}
		s.EventService = events.NewService(engineLogger, config.Events.AccessTokens)
		go s.EventService.Run()
		unauthGroup.POST("/stream/events/:token", s.EventService.PostNewEvent)
		unauthGroup.GET("/stream/events", s.EventService.AcceptGinWebsocket)
	} else {
		unauthGroup.GET("/stream/events", notImplemented)
		unauthGroup.POST("/stream/events/:token", notImplemented)
	}

	for level, endpoints := range routes {
		var group *gin.RouterGroup
		switch level {
		case AuthLevelNone:
			group = unauthGroup
		case AuthLevelUser:
			group = userGroup
		case AuthLevelAdmin:
			group = adminGroup
		default:
			panic("unknown auth level")
		}

		for _, route := range endpoints {
			switch route.method {
			case "GET":
				group.GET(route.endpoint, route.handler)
			case "POST":
				group.POST(route.endpoint, route.handler)
			case "PUT":
				group.PUT(route.endpoint, route.handler)
			case "DELETE":
				group.DELETE(route.endpoint, route.handler)
			}
		}
	}

	gamewebGroup := s.Engine.Group(config.GameWebPrefix)
	gamewebGroup.POST("/rank/mg3getrank.html", gameweb.PostGetRanks)
	gamewebGroup.GET("/text/:filename", gameweb.GetTextFile)
	gamewebGroup.POST("/reguser/reguser.html", gameweb.RegisterAccount)
	gamewebGroup.POST("/reguser/deluser.html", gameweb.DeleteAccount)
	gamewebGroup.POST("/reguser/chgpswd.html", gameweb.ChangePassword)

	return
}

func Success(c *gin.Context, data any) {
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
	Success(c, "Not implemented")
}
