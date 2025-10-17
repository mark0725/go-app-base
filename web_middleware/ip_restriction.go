package middleware

// IP Restriction
type IPRestrictionConfig struct {
	Allow   []string `json:"allow"`   // 允许的IP地址或CIDR, 例如: "192.168.1.1" 或 "192.168.0.0/16"
	Deny    []string `json:"deny"`    // 拒绝的IP地址或CIDR, 例如: "192.168.1.1" 或 "192.168.0.0/16"
	Message string   `json:"message"` // 拒绝请求时响应体中返回的信息
	Status  int      `json:"status"`  // 插件拒绝请求时使用的HTTP状态码
}
