package proxy

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"

	"github.com/nais/kolide-api-proxy/internal/cache"
	"github.com/nais/kolide-api-proxy/internal/kolide"
)

type response struct {
	Devices      []kolide.Device `json:"devices,omitempty"`
	LastModified string          `json:"last_modified,omitempty"`
}

func auth(username, password []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		u, p, ok := req.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(u), username) != 1 || subtle.ConstantTimeCompare([]byte(p), password) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Kolide API Proxy"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func apiDevices(c *cache.Cache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		devices := c.GetDevices()
		if len(devices) == 0 {
			w.Header().Set("Retry-After", "30")
			http.Error(w, http.StatusText(http.StatusServiceUnavailable)+": data not yet available, try again later", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{
			Devices:      devices,
			LastModified: c.LastModified().Format(http.TimeFormat),
		})
	})
}
