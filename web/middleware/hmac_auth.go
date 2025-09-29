package middleware

import "github.com/gin-gonic/gin"

// HmacAuthConfig defines the available settings for the HTTP-Signature authentication plugin.
type HmacAuthConfig struct {
	Algorithms          []string `json:"algorithms"`                      // HMAC algorithms to accept. Allowed: hmac-sha1 | hmac-sha256 | hmac-sha384 | hmac-sha512. Default: all four
	Anonymous           string   `json:"anonymous,omitempty"`             // Consumer UUID/username to use when authentication fails
	ClockSkew           int      `json:"clock_skew,omitempty"`            // Allowed clock skew (seconds) to prevent replay attacks. Default: 300
	EnforceHeaders      []string `json:"enforce_headers,omitempty"`       // Headers that MUST be included in the signature. Default: []
	HideCredentials     bool     `json:"hide_credentials,omitempty"`      // If true, omit credential when forwarding to upstream. Default: false
	Realm               string   `json:"realm,omitempty"`                 // Realm value returned in WWW-Authenticate on failure
	ValidateRequestBody bool     `json:"validate_request_body,omitempty"` // Enable request-body validation. Default: false
}

func (m *WebMiddlewareSecurity) HmacAuth(params map[string]any) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("auth_user_id", "xilin.gao")

	}
}
