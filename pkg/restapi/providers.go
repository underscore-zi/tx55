package restapi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"tx55/pkg/configurations"
)

// RequireAPIKey just requires that the X-API-TOKEN header is set to some value.
// The particular value does not matter as the intent is purely to determine if
// the call is has the ability to set custom headers on a request.
func RequireAPIKey(c *gin.Context) {
	if c.Request.Method == "POST" {
		apiKey := c.GetHeader("X-API-TOKEN")
		if apiKey == "" {
			Error(c, 401, "unauthorized")
			return
		}
	}
	c.Next()
}

func GameLoginRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user_id")
	if user == nil {
		Error(c, 401, "unauthorized")
		return
	} else {
		c.Next()
	}
}
func AdminLoginRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("admin_id")
	if user == nil {
		Error(c, 401, "unauthorized")
		return
	} else {
		c.Next()
	}
}

func ProvideContextVar(name string, val any) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(name, val)
		c.Next()
	}
}

func CORSMiddleware(config configurations.RestAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			// Not a CORS request, don't sent any CORS headers
		} else {
			// Check if the origin is in our list of allowed origins
			wasPresent := false
			for _, allowed := range config.AllowedCredentialOrigins {
				if origin == allowed {
					wasPresent = true
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
					c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-API-TOKEN, Authorization, accept, origin, Cache-Control, X-Requested-With")
					c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
					break
				}
			}
			if !wasPresent && len(config.AllowedOrigins) > 0 {
				if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
					c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					for _, allowed := range config.AllowedOrigins {
						if origin == allowed {
							c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
