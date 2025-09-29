package middleware

import (
	"github.com/gin-gonic/gin"
	base_utils "github.com/mark0725/go-app-base/utils"
	base_web "github.com/mark0725/go-app-base/web"
)

type BasicAuthConfig struct {
	Anonymous       string `json:"anonymous,omitempty" yaml:"anonymous,omitempty"`               // 可选字符串（Consumer UUID 或用户名），用作“anonymous”消费者。如果为空（默认null），认证失败时请求将返回4xx。该值只能为 Consumer 的 id 或 username，而非 custom_id。
	HideCredentials bool   `json:"hide_credentials,omitempty" yaml:"hide_credentials,omitempty"` // 可选布尔值。为 true 时，插件在代理前隐藏凭证（如Authorization头），不传递给上游服务。默认值为 false。
	Realm           string `json:"realm,omitempty" yaml:"realm,omitempty"`                       // 在认证失败时，插件发送 WWW-Authenticate 头，realm 属性为该值。默认值为 "service"。
}

func newBasicAuthConfig(params map[string]any) *BasicAuthConfig {
	conf := BasicAuthConfig{
		Anonymous:       "",
		HideCredentials: false,
	}

	base_utils.MapToStruct(params, &conf)

	return &conf
}
func (m *WebMiddlewareSecurity) BasicAuth(params map[string]any, r *gin.RouterGroup) gin.HandlerFunc {
	// conf := newBasicAuthConfig(params)
	return func(c *gin.Context) {
		if _, ok := c.Get(base_web.CtxKeyAuthenticatedConsumer); !ok {
			authedConsumer := base_web.AuthenticatedConsumer{
				Id:       "xilin.gao",
				Username: "高西林",
			}
			c.Set(base_web.CtxKeyAuthenticatedConsumer, &authedConsumer)
		}
		c.Next()

	}
}
