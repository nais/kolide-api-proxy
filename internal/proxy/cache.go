package proxy

import (
	"context"
	"time"

	"github.com/nais/kolide-api-proxy/internal/cache"
	"github.com/nais/kolide-api-proxy/internal/kolide"
	"github.com/sirupsen/logrus"
)

const syncInterval = time.Hour

// updateCache periodically fetches the list of devices from Kolide and updates the cache. This function blocks until
// the context is cancelled.
func updateCache(ctx context.Context, kac *kolide.Client, c *cache.Cache, log logrus.FieldLogger) error {
	for {
		func() {
			start := time.Now()
			devices, err := kac.GetDevices(ctx)
			if err != nil {
				log.WithError(err).Errorf("unable to update cache")
				return
			}

			c.SetDevices(devices)
			log.
				WithFields(logrus.Fields{
					"duration":    time.Since(start).String(),
					"duration_ms": time.Since(start).Milliseconds(),
					"devices":     len(devices),
				}).
				Info("cache updated")
		}()

		log.
			WithField("next_run_in", syncInterval.String()).
			Info("cache update run finished")

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(syncInterval):
		}
	}
}
