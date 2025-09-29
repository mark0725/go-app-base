package middleware

import (
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	base_utils "github.com/mark0725/go-app-base/utils"
	base_web "github.com/mark0725/go-app-base/web"
)

// SessionConfig 定义了会话的配置参数
type SessionConfig struct {
	AbsoluteTimeout         int      `json:"absolute_timeout"`          // 会话cookie的绝对超时时间（秒），超过即失效。默认 86400
	Audience                string   `json:"audience"`                  // 会话受众，例如 "my-application"。默认 "default"
	CookieDomain            string   `json:"cookie_domain"`             // cookie绑定的域名
	CookieHTTPOnly          bool     `json:"cookie_http_only"`          // 是否设置 HttpOnly，仅允许服务器访问cookie。默认 true
	CookieName              string   `json:"cookie_name"`               // cookie名称。默认 "session"
	CookiePath              string   `json:"cookie_path"`               // cookie所属的路径。默认 "/"
	CookieSameSite          string   `json:"cookie_same_site"`          // 跨站请求时cookie策略，允许 "Default"、"Lax"、"None"、"Strict"。默认 "Strict"
	CookieSecure            bool     `json:"cookie_secure"`             // 是否仅通过HTTPS发送cookie。默认 true
	HashSubject             bool     `json:"hash_subject"`              // store_metadata 启用时，是否对 subject 进行哈希。默认 false
	IdlingTimeout           int      `json:"idling_timeout"`            // 会话cookie的空闲超时时间（秒），默认 900
	LogoutMethods           []string `json:"logout_methods"`            // 插件支持响应的登出HTTP方法。Allowed: DELETE, GET, POST；默认 ["DELETE", "POST"]
	LogoutPostArg           string   `json:"logout_post_arg"`           // POST登出请求参数。默认 "session_logout"
	LogoutQueryArg          string   `json:"logout_query_arg"`          // query参数，登出时附带。默认 "session_logout"
	ReadBodyForLogout       bool     `json:"read_body_for_logout"`      // 登出时是否读取request body。默认 false
	Remember                bool     `json:"remember"`                  // 是否启用持久（长期）会话。默认 false
	RememberAbsoluteTimeout int      `json:"remember_absolute_timeout"` // 持久会话绝对超时时间（秒），默认 2592000
	RememberCookieName      string   `json:"remember_cookie_name"`      // 持久会话cookie名称。默认 "remember"
	RememberRollingTimeout  int      `json:"remember_rolling_timeout"`  // 持久会话rolling超时时间（秒），默认 604800
	RequestHeaders          []string `json:"request_headers"`           // 需要作为响应头包含下游的字段列表
	ResponseHeaders         []string `json:"response_headers"`          // 需要作为响应头包含下游的字段列表
	RollingTimeout          int      `json:"rolling_timeout"`           // 普通会话rolling超时时间（秒），默认 3600
	Secret                  string   `json:"secret"`                    // 用于HMAC签名的密钥，需安全加密。默认 "M5N0CyFrAZwpB7F72PpY3J4S5n3KL77fs6xMuY8b7SMc"
	StaleTTL                int      `json:"stale_ttl"`                 // 会话失效后旧cookie丢弃时间（秒）。默认 10
	Storage                 string   `json:"storage"`                   // 会话存储方式，可选 "cookie" 或 "kong"。默认 "cookie"
	StoreMetadata           bool     `json:"store_metadata"`            // 是否存储元数据，例如收集特定audience、subject的会话信息。默认 false
}

func newSessionConfig(params map[string]any) *SessionConfig {
	conf := SessionConfig{
		AbsoluteTimeout:         0,
		Audience:                "default",
		CookieDomain:            "",
		CookieHTTPOnly:          false,
		CookieName:              "session",
		CookiePath:              "/",
		CookieSameSite:          "Strict",
		CookieSecure:            true,
		HashSubject:             false,
		IdlingTimeout:           900,
		LogoutMethods:           []string{"DELETE", "POST"},
		LogoutPostArg:           "session_logout",
		LogoutQueryArg:          "session_logout",
		ReadBodyForLogout:       false,
		Remember:                false,
		RememberAbsoluteTimeout: 2592000,
		RememberCookieName:      "remember",
		RememberRollingTimeout:  604800,
		RequestHeaders:          []string{},
		ResponseHeaders:         []string{},
		RollingTimeout:          3600,
		Secret:                  "M5N0CyFrAZwpB7F72PpY3J4S5n3KL77fs6xMuY8b7SMc",
		StaleTTL:                10,
		Storage:                 "cookie",
		StoreMetadata:           false,
	}

	base_utils.MapToStruct(params, &conf)

	return &conf
}
func (m *WebMiddlewareSecurity) Session(params map[string]any, r *gin.RouterGroup) gin.HandlerFunc {
	conf := newSessionConfig(params)
	store := cookie.NewStore([]byte(conf.Secret))
	r.Use(sessions.Sessions(conf.CookieName, store))

	return func(c *gin.Context) {
		session := sessions.Default(c)
		if v := session.Get(base_web.CtxKeyAuthenticatedConsumer); v != nil {
			authedConsumer := base_web.AuthenticatedConsumer{}
			if str, ok := v.(string); ok {
				if err := json.Unmarshal([]byte(str), &authedConsumer); err == nil {
					c.Set(base_web.CtxKeyAuthenticatedConsumer, &authedConsumer)
				}
			}
		}
		c.Next()
		if authedConsumer, ok := c.Get(base_web.CtxKeyAuthenticatedConsumer); ok {
			if data, err := json.Marshal(authedConsumer); err == nil {
				session.Set(base_web.CtxKeyAuthenticatedConsumer, data) // 设置 session
				session.Save()
			}
		}

	}
}
