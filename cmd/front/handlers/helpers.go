package handlers

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIp(r *http.Request) string {
	//Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return stripPort(strings.TrimSpace(ips[0]))
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return stripPort(xri)
	}

	//Last resort - take RemoteAddr directly
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr //IPv6 or no port
	}
	return stripPort(ip)
}

// Removes the port part from an ip address string
func stripPort(addr string) string {
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	//No port or malformed address
	return addr
}
