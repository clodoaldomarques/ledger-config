package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clodoaldomarques/core-sdk/pkg/otel"
	"github.com/clodoaldomarques/core-sdk/pkg/zap/logger"
	"github.com/clodoaldomarques/ledger-config/config"
	"github.com/clodoaldomarques/ledger-config/internal/infra/rest/server"
)

func main() {
	ctx := context.Background()
	otel.Start(ctx)

	s := server.New()
	go func() {
		c := config.New()
		err := s.Start(c.AppPort)
		if err != http.ErrServerClosed {
			logger.Fatal(ctx, err.Error(), logger.Fields{})
		}
	}()

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

	if err := otel.Shutdown(ctx); err != nil {
		fmt.Printf("erro no shutdown do opentelemetry: %v\n", err)
	}

	fmt.Println("Graceful shutdown concluído")
}
