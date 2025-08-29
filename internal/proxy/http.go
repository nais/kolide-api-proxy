package proxy

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/nais/kolide-api-proxy/internal/cache"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func runHTTPServer(ctx context.Context, cache *cache.Cache, listenAddr string, proxyApiToken string, log logrus.FieldLogger) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(http.ResponseWriter, *http.Request) {})
	mux.Handle("GET /api/devices", auth(proxyApiToken, apiDevices(cache)))

	srv := &http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Infof("HTTP server shutting down...")
		if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.WithError(err).Errorf("HTTP server shutdown failed")
			return err
		}
		return nil
	})

	eg.Go(func() error {
		log.WithField("addr", listenAddr).Infof("HTTP server accepting requests")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Errorf("unexpected error from HTTP server")
			return err
		}
		log.Infof("HTTP server finished, terminating...")
		return nil
	})

	return eg.Wait()
}
