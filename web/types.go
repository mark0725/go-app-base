package web

// ListenConfig 用于存储解析后的listen配置
type ListenConfig struct {
	Address       string
	Port          int
	SSL           bool
	HTTP2         bool
	ProxyProtocol bool
	Deferred      bool
	Bind          bool
	ReusePort     bool
	Backlog       int
	IPv6Only      *bool
	SOKeepAlive   *ListenKeepAliveConfig
}

// KeepAliveConfig 用于存储TCP keepalive设置
type ListenKeepAliveConfig struct {
	On        bool
	KeepIdle  int
	KeepIntvl int
	KeepCnt   int
}

const CtxKeyAuthenticatedConsumer = "authenticated_consumer"

type AuthenticatedConsumer struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	OrgId    string `json:"org_id"`
}

const CtxKeyAuthenticatedCredential = "authenticated_credential"

type AuthenticatedCredential struct {
	Id         string `json:"id"`
	Identifier string `json:"identifier"`
}

const CtxKeyAuthenticatedGroups = "authenticated_groups"

// authenticated_groups
type AuthenticatedGroups []string
