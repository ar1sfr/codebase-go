package cors

import (
	"strings"

	"codebase-go/internal/config"

	"github.com/gin-gonic/gin"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// handle allowed origins
		if len(cfg.CORS.AllowOrigins) > 0 {
			if cfg.CORS.AllowOrigins[0] == "*" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				for _, allowedOrigin := range cfg.CORS.AllowOrigins {
					if origin == allowedOrigin {
						c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}
		}

		if cfg.CORS.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if len(cfg.CORS.AllowHeaders) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.CORS.AllowHeaders, ","))
		}

		if len(cfg.CORS.AllowMethods) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.CORS.AllowMethods, ","))
		}

		if len(cfg.CORS.ExposeHeaders) > 0 {
			c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(cfg.CORS.ExposeHeaders, ","))
		}

		if cfg.CORS.MaxAge > 0 {
			c.Writer.Header().Set("Access-Control-Max-Age", string(cfg.CORS.MaxAge))
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
