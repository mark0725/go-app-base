package middleware

import (
	"github.com/gin-gonic/gin"
	base_utils "github.com/mark0725/go-app-base/utils"
	base_web "github.com/mark0725/go-app-base/web"
)

// LDAPAuthConfig defines the configuration parameters for the LDAP authentication plugin.
type LDAPAuthConfig struct {
	Anonymous            string   `json:"anonymous,omitempty"`              // Optional consumer UUID or username used when authentication fails (empty = 4xx).
	Attribute            string   `json:"attribute"`                        // LDAP attribute to match the user (e.g., "cn").
	BaseDN               string   `json:"base_dn"`                          // Base distinguished name to start the user search (e.g., "dc=example,dc=com").
	BindDN               string   `json:"bind_dn,omitempty"`                // DN with privileges to perform the search.
	CacheTTL             int      `json:"cache_ttl,omitempty"`              // Cache expiry time in seconds (default 60).
	ConsumerBy           []string `json:"consumer_by,omitempty"`            // Fields used to map the consumer: "custom_id", "username", or both.
	ConsumerOptional     bool     `json:"consumer_optional,omitempty"`      // Skip consumer mapping if true (default false).
	GroupBaseDN          string   `json:"group_base_dn,omitempty"`          // DN where group searches begin.
	GroupMemberAttribute string   `json:"group_member_attribute,omitempty"` // Attribute holding group members (default "memberOf").
	GroupNameAttribute   string   `json:"group_name_attribute,omitempty"`   // Attribute holding group name (e.g., "cn").
	GroupsRequired       []string `json:"groups_required,omitempty"`        // LDAP groups required for successful authorization.
	HeaderType           string   `json:"header_type,omitempty"`            // Token type placed in Authorization header (default "ldap").
	HideCredentials      bool     `json:"hide_credentials,omitempty"`       // Remove credentials before proxying if true (default false).
	Keepalive            int      `json:"keepalive,omitempty"`              // Idle LDAP connection lifetime in ms (default 60000).
	LdapHost             string   `json:"ldap_host"`                        // LDAP server hostname.
	LdapPassword         string   `json:"ldap_password,omitempty"`          // Password to bind to the LDAP server.
	LdapPort             int      `json:"ldap_port,omitempty"`              // LDAP server port (default 389, or 636 for LDAPS).
	Ldaps                bool     `json:"ldaps,omitempty"`                  // Use LDAPS (TLS/SSL) if true (default false).
	LogSearchResults     bool     `json:"log_search_results,omitempty"`     // Log full search results for debugging (default false).
	Realm                string   `json:"realm,omitempty"`                  // Realm returned in WWW-Authenticate header on failure.
	StartTLS             bool     `json:"start_tls,omitempty"`              // Use StartTLS over LDAP connection (default false).
	Timeout              int      `json:"timeout,omitempty"`                // LDAP connection timeout in ms (default 10000).
	VerifyLdapHost       bool     `json:"verify_ldap_host,omitempty"`       // Verify LDAP server certificate if true (default false).
}

func newLDAPAuthConfig(params map[string]any) *LDAPAuthConfig {
	conf := LDAPAuthConfig{
		Anonymous:       "",
		HideCredentials: false,
	}

	base_utils.MapToStruct(params, &conf)

	return &conf
}
func (m *WebMiddlewareSecurity) LdapAuth(params map[string]any, r *gin.RouterGroup) gin.HandlerFunc {
	// conf := newLDAPAuthConfig(params)
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
