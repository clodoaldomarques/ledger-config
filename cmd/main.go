package main

import (
	"context"
	"net/http"

	"github.com/clodoaldomarques/core-sdk/pkg/logger"
	"github.com/clodoaldomarques/ledger-config/configs"
	"github.com/clodoaldomarques/ledger-config/internal/infra/rest/server"
)

func main() {
	c := configs.New(configs.WithAppPort(5000))
	err := server.New().Start(c.AppPort)
	if err != http.ErrServerClosed {
		logger.Fatal(context.Background(), err.Error(), logger.Fields{})
	}
}
