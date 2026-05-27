package main

import (
	"context"
	"net/http"

	"github.com/clodoaldomarques/core-sdk/pkg/logger"
	"github.com/clodoaldomarques/ledger-config/config"
	"github.com/clodoaldomarques/ledger-config/internal/infra/rest/server"
)

func main() {
	c := config.New(config.WithAppPort(5000))
	err := server.New().Start(c.AppPort)
	if err != http.ErrServerClosed {
		logger.Fatal(context.Background(), err.Error(), logger.Fields{})
	}
}
