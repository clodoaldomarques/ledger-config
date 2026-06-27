package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clodoaldomarques/core-sdk/pkg/logger"
	"github.com/clodoaldomarques/core-sdk/pkg/opentelemetry"
	"github.com/clodoaldomarques/ledger-config/config"
	"github.com/clodoaldomarques/ledger-config/internal/infra/rest/server"
)

func main() {

	s := server.New()
	go func() {
		c := config.New(config.WithAppPort(5000))
		err := s.Start(c.AppPort)
		if err != http.ErrServerClosed {
			logger.Fatal(context.Background(), err.Error(), logger.Fields{})
		}
	}()

	opentelemetry.Start(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\nRecebido sinal de desligamento, iniciando graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		fmt.Printf("Erro no shutdown do servidor HTTP: %v\n", err)
	} else {
		fmt.Println("Servidor HTTP finalizado com sucesso")
	}

	if err := opentelemetry.Shutdown(ctx); err != nil {
		fmt.Printf("erro no shutdown do opentelemetry: %v\n", err)
	}

	fmt.Println("Graceful shutdown concluído")
}
