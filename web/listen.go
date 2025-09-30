package web

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseListenConfig 解析给定的proxy_listen配置字符串
func ParseListenConfig(proxyListen string) ([]*ListenConfig, error) {
	var configs []*ListenConfig
	configStrings := strings.Split(proxyListen, ",")

	for _, configStr := range configStrings {
		configStr = strings.TrimSpace(configStr)

		var config ListenConfig

		parts := strings.Fields(configStr)
		if len(parts) < 1 {
			return nil, fmt.Errorf("invalid listen config: %s", configStr)
		}
		addrPort := strings.Split(parts[0], ":")

		if len(addrPort) != 2 {
			return nil, fmt.Errorf("invalid address:port format")
		}

		config.Address = addrPort[0]
		port, err := strconv.Atoi(addrPort[1])
		if err != nil {
			return nil, fmt.Errorf("invalid port number")
		}
		config.Port = port

		for _, part := range parts[1:] {
			switch {
			case part == "ssl":
				config.SSL = true
			case part == "http2":
				config.HTTP2 = true
			case part == "proxy_protocol":
				config.ProxyProtocol = true
			case part == "deferred":
				config.Deferred = true
			case part == "bind":
				config.Bind = true
			case part == "reuseport":
				config.ReusePort = true
			case strings.HasPrefix(part, "backlog="):
				config.Backlog, err = strconv.Atoi(strings.TrimPrefix(part, "backlog="))
				if err != nil {
					return nil, fmt.Errorf("invalid backlog value")
				}
			case strings.HasPrefix(part, "ipv6only="):
				value := strings.TrimPrefix(part, "ipv6only=")
				config.IPv6Only = parseOnOff(value)
			case strings.HasPrefix(part, "so_keepalive="):
				config.SOKeepAlive = parseKeepAlive(strings.TrimPrefix(part, "so_keepalive="))
				if config.SOKeepAlive == nil {
					return nil, fmt.Errorf("invalid keepalive configuration")
				}
			default:
				return nil, fmt.Errorf("unknown parameter: %s", part)
			}
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

func parseOnOff(value string) *bool {
	var result bool
	switch value {
	case "on":
		result = true
	case "off":
		result = false
	default:
		return nil
	}
	return &result
}

func parseKeepAlive(value string) *ListenKeepAliveConfig {
	if value == "on" {
		return &ListenKeepAliveConfig{On: true}
	} else if value == "off" {
		return &ListenKeepAliveConfig{On: false}
	} else {
		params := strings.Split(value, ":")
		if len(params) != 3 {
			return nil
		}
		keepIdle, err1 := strconv.Atoi(params[0])
		keepIntvl, err2 := strconv.Atoi(params[1])
		keepCnt, err3 := strconv.Atoi(params[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return nil
		}
		return &ListenKeepAliveConfig{KeepIdle: keepIdle, KeepIntvl: keepIntvl, KeepCnt: keepCnt}
	}

}
