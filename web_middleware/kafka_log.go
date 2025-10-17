package middleware

// Config is the root configuration object.
type KafkaLogConfig struct {
	Authentication   Authentication    `json:"authentication"`    // SASL authentication block
	BootstrapServers []BootstrapServer `json:"bootstrap_servers"` // List of bootstrap brokers
}

// Authentication holds SASL-related settings.
type Authentication struct {
	Mechanism string `json:"mechanism"` // PLAIN | SCRAM-SHA-256 | SCRAM-SHA-512
	Password  string `json:"password"`  // Password (encrypted / referenceable)
	Strategy  string `json:"strategy"`  // Authentication strategy, fixed value "sasl"
	TokenAuth bool   `json:"tokenauth"` // Delegation-token authentication switch
	User      string `json:"user"`      // Username (encrypted / referenceable)
}

// BootstrapServer represents a single broker endpoint.
type BootstrapServer struct {
	Host string `json:"host"` // Host name (e.g. example.com)
	Port int    `json:"port"` // Port 0 – 65535
}
