package proxy

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func newLogger(level string) (logrus.FieldLogger, error) {
	log := logrus.StandardLogger()
	log.SetFormatter(&logrus.JSONFormatter{})

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}

	log.SetLevel(parsedLevel)

	return log, nil
}
