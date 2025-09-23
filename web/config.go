package web

type WebConfig struct {
	Oauth   WebOauthConfig           `json:"oauth"`
	Certs   map[string]WebCertConfig `json:"certs"`
	Servers []WebServer              `json:"servers"`
}

type WebServer struct {
	Name        string                `json:"name"`
	Listen      string                `json:"listen"`
	Cert        string                `json:"cert"`
	Log         WebLogConfig          `json:"log"`
	Session     WebSessionConfig      `json:"session"`
	Configs     []WebServerConfig     `json:"configs"`
	Middlewares []WebMiddlewareConfig `json:"middlewares"`
	Endpoints   []WebEndpoint         `json:"endpoints"`
	TemplateDir string                `json:"template_dir"`
	StaticDir   string                `json:"static_dir"`
}
type WebOauthConfig struct {
	CodeExpireIn  int `json:"code_expire_in"`
	TokenExpireIn int `json:"token_expire_in"`
	AuthExpireIn  int `json:"auth_expire_in"`
}

type WebServerConfig struct {
	Name   string         `json:"name"`
	Module string         `json:"module"`
	Params map[string]any `json:"params,omitempty"`
}

type WebMiddlewareConfig struct {
	Name   string         `json:"name"`
	Module string         `json:"module"`
	Params map[string]any `json:"params,omitempty"`
}

type WebEndpoint struct {
	Path        string                `json:"path"`
	Module      string                `json:"module"`
	Group       string                `json:"group"`
	Middlewares []WebMiddlewareConfig `json:"middlewares,omitempty"`
}

type WebSessionConfig struct {
	Name         string `json:"name"`
	Duration     int    `json:"duration"`
	SessionStore string `json:"session_store"`
	CookieSecret string `json:"cookie_secret"`
}

type WebLogConfig struct {
	Level     string `json:"level"`
	AccessLog string `json:"access_log"`
	ErrorLog  string `json:"error_log"`
}

type WebCertConfig struct {
	Name     string `json:"name"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}
