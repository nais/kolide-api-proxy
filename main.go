package main

import (
	"context"

	"github.com/nais/kolide-api-proxy/internal/proxy"
)

func main() {
	proxy.Run(context.Background())
}
