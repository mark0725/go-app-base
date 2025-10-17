package middleware

// OAuth2Config represents the configurable parameters for the
// OAuth 2.0  plugin.
type OAuth2Config struct {
	Anonymous           string            `json:"anonymous,omitempty"`             // optional consumer UUID or username used as “anonymous” consumer when auth fails (null to disable)
	AuthorizationValue  string            `json:"authorization_value"`             // required, value of Authorization header for the  request (usually “Basic <base-64>”)
	ConsumerBy          string            `json:"consumer_by,omitempty"`           // associates OAuth2 username or client_id with consumer ("username" | "client_id"), default "username"
	CustomClaimsForward []string          `json:"custom_claims_forward,omitempty"` // list of custom claims to forward to upstream as X-Credential-{claim}
	CustomHeaders       map[string]string `json:"custom__headers,omitempty"`       // additional custom headers to include in the  request
	HideCredentials     bool              `json:"hide_credentials,omitempty"`      // if true, removes credential before proxying to upstream, default false
	IntrospectRequest   bool              `json:"introspect_request,omitempty"`    // if true, forwards X-Request-Path and X-Request-Http-Method to the  endpoint
	URL                 string            `json:"url"`                             // required URL of the  endpoint (e.g., https://example.com/introspect)
	KeepAlive           int               `json:"keepalive,omitempty"`             // idle connection lifetime in ms before close, default 60000
	RunOnPreflight      bool              `json:"run_on_preflight,omitempty"`      // if false, skips auth on OPTIONS preflight requests, default true
	Timeout             int               `json:"timeout,omitempty"`               // timeout in ms when sending data to upstream server, default 10000
	TokenTypeHint       string            `json:"token_type_hint,omitempty"`       // token_type_hint value included in  requests
	TTL                 float64           `json:"ttl,omitempty"`                   // TTL (seconds) for caching  response, 0 disables, default 30
}
