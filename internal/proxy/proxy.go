package proxy

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nais/kolide-api-proxy/internal/cache"
	"github.com/nais/kolide-api-proxy/internal/kolide"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	exitCodeSuccess = iota
	exitCodeEnvFileError
	exitCodeConfigError
	exitCodeLoggerError
	exitCodeRunError
)

func Run(ctx context.Context) {
	log := logrus.StandardLogger()

	if err := loadEnvFile(log); err != nil {
		log.WithError(err).Errorf("error loading .env file")
		os.Exit(exitCodeEnvFileError)
	}

	cfg, err := newConfig(ctx)
	if err != nil {
		log.WithError(err).Errorf("error when loading config")
		os.Exit(exitCodeConfigError)
	}

	appLogger, err := newLogger(cfg.LogLevel)
	if err != nil {
		log.WithError(err).Errorf("creating application logger")
		os.Exit(exitCodeLoggerError)
	}

	if err := run(ctx, cfg, appLogger); err != nil {
		appLogger.WithError(err).Errorf("error in run()")
		os.Exit(exitCodeRunError)
	}

	os.Exit(exitCodeSuccess)
}

func run(ctx context.Context, cfg *Config, log logrus.FieldLogger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx, shutdown := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer shutdown()

	c := cache.New()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		kac := kolide.NewClient(cfg.KolideApiToken, log.WithField("component", "kolideapiclient"))
		return updateCache(ctx, kac, c, log.WithField("component", "cache_updater"))
	})

	eg.Go(func() error {
		return runHTTPServer(ctx, c, cfg.ListenAddress, cfg.ProxyApiToken, log.WithField("component", "http"))
	})

	<-ctx.Done()
	shutdown()

	ch := make(chan error)
	go func() { ch <- eg.Wait() }()

	select {
	case <-time.After(10 * time.Second):
		log.Warn("timed out waiting for graceful shutdown")
	case err := <-ch:
		return err
	}

	return nil
}
