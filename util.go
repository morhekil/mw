package mw

import (
	"net/http"
	"strings"
)

// extract real IP from various possible proxy headers
func realIP(req *http.Request) string {
	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
	}
	if addr == "" {
		addr = strings.SplitN(req.RemoteAddr, ":", 1)[0]
	}

	return addr
}
