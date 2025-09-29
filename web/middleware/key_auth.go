package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	base_utils "github.com/mark0725/go-app-base/utils"
	base_web "github.com/mark0725/go-app-base/web"
)

// ConfigObject represents the configuration schema described above.
type KeyAuthConfig struct {
	Anonymous       string          `json:"anonymous,omitempty"`        // optional consumer UUID or username used when authentication fails
	HideCredentials bool            `json:"hide_credentials,omitempty"` // if true, strips the credential from the request (default: false)
	IdentityRealms  []IdentityRealm `json:"identity_realms,omitempty"`  // list of Konnect Identity Realms used to source a consumer
	KeyInBody       bool            `json:"key_in_body,omitempty"`      // enables reading the credential from request body (default: false)
	KeyInHeader     bool            `json:"key_in_header,omitempty"`    // enables reading the credential from request header (default: true)
	KeyInQuery      bool            `json:"key_in_query,omitempty"`     // enables reading the credential from query parameter (default: true)
	KeyNames        []string        `json:"key_names,omitempty"`        // header/query/body parameter names to look for (default: ["apikey"])
	Realm           string          `json:"realm,omitempty"`            // realm value returned in WWW-Authenticate header on auth failure
	RunOnPreflight  bool            `json:"run_on_preflight,omitempty"` // if true, authenticates OPTIONS preflight requests (default: true)
}

// IdentityRealm defines an individual Konnect Identity Realm entry.
type IdentityRealm struct {
	ID     string `json:"id,omitempty"`     // UUID of the identity realm
	Region string `json:"region,omitempty"` // region identifier
	Scope  string `json:"scope,omitempty"`  // allowed values: "cp" or "realm"
}

func newKeyAuthConfig(params map[string]any) *KeyAuthConfig {
	conf := KeyAuthConfig{
		Anonymous:       "",
		HideCredentials: false,
		IdentityRealms:  []IdentityRealm{},
		KeyInBody:       false,
		KeyInHeader:     false,
		KeyInQuery:      false,
		KeyNames:        []string{},
		Realm:           "",
		RunOnPreflight:  false,
	}

	base_utils.MapToStruct(params, &conf)

	return &conf
}
func (m *WebMiddlewareSecurity) KeyAuth(params map[string]any, r *gin.RouterGroup) gin.HandlerFunc {
	// conf := newKeyAuthConfig(params)
	return func(c *gin.Context) {
		// apiKey, err := extractAPIKeyFromHeader(c)
		// if err != nil {
		// 	c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
		// 	return
		// }

		// if apiInfo, err := repo.GetByApiKey(apiKey); err != nil {
		// 	c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
		// 	return
		// } else {
		// 	c.Set("api_info", apiInfo)
		// }

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

func extractAPIKeyFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid Authorization header format")
	}

	return strings.TrimSpace(parts[1]), nil
}
