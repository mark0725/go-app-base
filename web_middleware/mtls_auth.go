package middleware

// Config defines the plugin configuration.
type MTLSConfig struct {
	AllowPartialChain    bool     `json:"allow_partial_chain,omitempty"`    // Allow certificate verification with only an intermediate certificate. Default: false
	Anonymous            string   `json:"anonymous,omitempty"`              // Consumer UUID or username to use as “anonymous” if authentication fails
	AuthenticatedGroupBy string   `json:"authenticated_group_by,omitempty"` // Certificate property to use as the authenticated group (CN or DN). Default: CN
	CaCertificates       []string `json:"ca_certificates"`                  // List of CA certificate IDs to validate a client certificate (required)
	CacheTTL             int64    `json:"cache_ttl,omitempty"`              // Cache expiry time in seconds. Default: 60
	CertCacheTTL         int64    `json:"cert_cache_ttl,omitempty"`         // Seconds between refreshes of the revocation-status cache. Default: 60000
	ConsumerBy           []string `json:"consumer_by,omitempty"`            // Fields used to auto-match certificate subject to a consumer. Default: ["custom_id","username"]
	DefaultConsumer      string   `json:"default_consumer,omitempty"`       // Consumer to use when a trusted certificate is presented but no consumer matches
	HTTPProxyHost        string   `json:"http_proxy_host,omitempty"`        // HTTP proxy host for OCSP/CRL requests
	HTTPProxyPort        int      `json:"http_proxy_port,omitempty"`        // HTTP proxy port (0-65535)
	HTTPTimeout          int64    `json:"http_timeout,omitempty"`           // HTTP timeout (ms) when talking to OCSP/CRL servers. Default: 30000
	HTTPSProxyHost       string   `json:"https_proxy_host,omitempty"`       // HTTPS proxy host for OCSP/CRL requests
	HTTPSProxyPort       int      `json:"https_proxy_port,omitempty"`       // HTTPS proxy port (0-65535)
	RevocationCheckMode  string   `json:"revocation_check_mode,omitempty"`  // Certificate revocation check mode: IGNORE_CA_ERROR, SKIP, or STRICT. Default: IGNORE_CA_ERROR
	SendCADN             bool     `json:"send_ca_dn,omitempty"`             // Send CA distinguished names in the TLS handshake. Default: false
	SkipConsumerLookup   bool     `json:"skip_consumer_lookup,omitempty"`   // Skip consumer lookup once the certificate is trusted. Default: false
}
