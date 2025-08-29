package proxy

import (
	"encoding/json"
	"net/http"

	"github.com/nais/kolide-api-proxy/internal/cache"
	"github.com/nais/kolide-api-proxy/internal/kolide"
)

type response struct {
	Devices      []kolide.Device `json:"devices,omitempty"`
	LastModified string          `json:"last_modified,omitempty"`
}

func auth(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "Bearer "+token {
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
			http.Error(w, "data not available, try again later", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{
			Devices:      devices,
			LastModified: c.LastModified().Format(http.TimeFormat),
		})
	})
}
