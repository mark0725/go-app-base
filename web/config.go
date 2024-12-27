package web

type WebConfig struct {
	Oauth   WebOauthConfig `json:"oauth"`
	Servers []WebServer    `json:"servers"`
}

type WebServer struct {
	Name          string        `json:"name"`
	Listen        string        `json:"listen"`
	SSLCert       string        `json:"ssl_cert"`
	SSLCertKey    string        `json:"ssl_cert_key"`
	StaticDir     string        `json:"static_dir"`
	TemplateDir   string        `json:"template_dir"`
	CookieSecret  string        `json:"cookie_secret"`
	SessionExpire int           `json:"session_expire"`
	SessionName   string        `json:"session_name"`
	Endpoints     []WebEndpoint `json:"endpoints"`
	AccessLog     string        `json:"access_log"`
	ErrorLog      string        `json:"error_log"`
}
type WebOauthConfig struct {
	CodeExpireIn  int `json:"code_expire_in"`
	TokenExpireIn int `json:"token_expire_in"`
	AuthExpireIn  int `json:"auth_expire_in"`
}

type WebEndpoint struct {
	Path   string `json:"path"`
	Module string `json:"module"`
	Group  string `json:"group"`
}
