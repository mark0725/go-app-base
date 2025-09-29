package middleware

import (
	"github.com/gin-gonic/gin"
	base_web "github.com/mark0725/go-app-base/web"
)

func init() {
	base_web.EndpointMiddlewareRegister("base", BaseWebMiddleware)
}

func BaseWebMiddleware(name string, params map[string]any, r *gin.RouterGroup) {
	switch name {
	//Authentication
	case "basic-auth":
		r.Use(g_WebMiddlewareSecurity.BasicAuth(params, r))
	case "key-auth":
		r.Use(g_WebMiddlewareSecurity.KeyAuth(params, r))
	case "hmac-auth":
		r.Use(g_WebMiddlewareSecurity.BasicAuth(params, r))
	case "jwt":
		r.Use(g_WebMiddlewareSecurity.Jwt(params, r))
	case "oauth2":
	case "openid-connect":
	case "ldap-auth":
	case "mtls-auth":
	case "session":
		r.Use(g_WebMiddlewareSecurity.Session(params, r))
	//Security
	case "cors":
		r.Use(g_WebMiddlewareSecurity.CORS(params, r))
	case "ip-restriction":
	//Traffic Control
	case "rate-limiting":
	case "request-termination":
		//Transformations
		//Logging
	case "http-log":
	case "file-log":
	case "kafka	-log":
		//Analytics & Monitoring
	case "opentelemetry":
	case "prometheus":

	}
}

type WebMiddlewareSecurity struct{}

var g_WebMiddlewareSecurity WebMiddlewareSecurity = WebMiddlewareSecurity{}
