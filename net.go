package tools

import (
	"net/http"
	"strings"
)

func GetIPFromRequest(r *http.Request) string {
	list := strings.Split(r.RemoteAddr, ":")
	if len(list) != 2 {
		return ""
	}
	return list[0]
}
