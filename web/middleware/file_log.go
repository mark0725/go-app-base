package middleware

// Config holds all available settings for the plugin.
type FileLogConfig struct {
	CustomFieldsByLua map[string]string `json:"custom_fields_by_lua,omitempty"` // Lua code as a key-value map, additional properties allowed
	Path              string            `json:"path"`                           // The file path of the output log file; created if it doesn’t exist
	Reopen            bool              `json:"reopen,omitempty"`               // Close and reopen the log file on every request (default: false)
}
