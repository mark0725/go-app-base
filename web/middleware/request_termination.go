package middleware

// Request Termination
type RequestTerminationConfig struct {
	Body        string `json:"body"`         // The raw response body to send. Mutually exclusive with the Message field.
	ContentType string `json:"content_type"` // Content type of the raw response configured with Body.
	Echo        bool   `json:"echo"`         // When set, echoes a copy of the request back to the client. Main use: debugging. Can be combined with Trigger to debug on live systems. Default: false.
	Message     string `json:"message"`      // The message to send, if using the default response generator.
	StatusCode  int    `json:"status_code"`  // The response code to send. Integer between 100 and 599. Default: 503.
	Trigger     string `json:"trigger"`      // A string representing an HTTP header name.
}
