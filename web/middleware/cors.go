package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	base_utils "github.com/mark0725/go-app-base/utils"
)

// CORSConfig represents the CORS configuration object.
type CORSConfig struct {
	AllowOriginAbsent bool     `json:"allow_origin_absent"` // Skip CORS response headers when request's Origin header is empty (default: true)
	Credentials       bool     `json:"credentials"`         // Whether to send Access-Control-Allow-Credentials header with "true" (default: false)
	ExposedHeaders    []string `json:"exposed_headers"`     // Custom headers to expose via Access-Control-Expose-Headers
	Headers           []string `json:"headers"`             // Allowed headers for Access-Control-Allow-Headers
	MaxAge            int      `json:"max_age"`             // Max age (seconds) to cache preflight request results
	Methods           []string `json:"methods"`             // Allowed HTTP methods for Access-Control-Allow-Methods
	Origins           []string `json:"origins"`             // Allowed domains for Access-Control-Allow-Origin (use "*" to allow all)
	PreflightContinue bool     `json:"preflight_continue"`  // Proxy OPTIONS preflight request to upstream service (default: false)
	PrivateNetwork    bool     `json:"private_network"`     // Whether to send Access-Control-Allow-Private-Network header (default: false)
}

func newCORSConfig(params map[string]any) *CORSConfig {
	conf := CORSConfig{}

	base_utils.MapToStruct(params, &conf)

	return &conf
}

func (m *WebMiddlewareSecurity) CORS(params map[string]any, r *gin.RouterGroup) gin.HandlerFunc {
	conf := newCORSConfig(params)

	return func(c *gin.Context) {
		// if conf.PreflightContinue {
		// 	c.AbortWithStatus(200)
		// }
		c.Next()
		if conf.Credentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if conf.AllowOriginAbsent {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			c.Header("Access-Control-Allow-Origin", strings.Join(conf.Origins, ", "))
		}
		if conf.ExposedHeaders != nil {
			c.Header("Access-Control-Expose-Headers", strings.Join(conf.ExposedHeaders, ", "))
		}
		if conf.Methods != nil {
			c.Header("Access-Control-Allow-Methods", strings.Join(conf.Methods, ", "))
		}
		if conf.Headers != nil {
			c.Header("Access-Control-Allow-Headers", strings.Join(conf.Headers, ", "))
		}
		if conf.PrivateNetwork {
			c.Header("Access-Control-Allow-Private-Network", "true")
		}

	}
}
