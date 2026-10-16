package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	otelecho "github.com/clodoaldomarques/core-sdk/pkg/otel/echo"
	"github.com/clodoaldomarques/core-sdk/pkg/otel/tracer"
	"github.com/clodoaldomarques/core-sdk/pkg/zap/logger"

	"github.com/clodoaldomarques/ledger-config/internal/infra/db/dynamodb"
	"github.com/clodoaldomarques/ledger-config/internal/infra/message"
	"github.com/clodoaldomarques/ledger-config/internal/infra/rest/ledger"
	"github.com/clodoaldomarques/ledger-config/internal/infra/rest/shared"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type Server struct {
	http *echo.Echo
}

func New() *Server {
	s := &Server{
		http: echo.New(),
	}
	s.routes()
	return s
}

func (s Server) routes() {
	s.http.Validator = &CustomValidator{validator: validator.New()}

	// tracer interceptor
	s.http.Use(tracer.Interceptor())

	// logger interceptor
	s.http.Use(logger.Interceptor())

	// custom metrics interceptor
	s.http.Use(otelecho.Interceptor())

	// Health
	s.http.GET("/health/live", s.Liveness)
	s.http.GET("/health/ready", s.Readiness)

	// Ledger handler
	s.http.POST("/v1/ledger/config", ledger.CreateConfig)
	s.http.PATCH("/v1/ledger/config/:config_id", ledger.UpdateConfig)
	s.http.DELETE("/v1/ledger/config/:config_id/disable", ledger.DisableConfig)

	s.http.GET("/v1/ledger/config/:event_type_id/:program_id", ledger.FindLedgerConfig)
	s.http.GET("/v1/ledger/configs", ledger.FindAllLedgerConfig)

	s.http.POST("/v1/ledger/org/activate", ledger.ActivateTenant)

}

func (s Server) Start(port int) error {
	return s.http.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *Server) Liveness(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (s *Server) Readiness(c echo.Context) error {
	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		500*time.Millisecond,
	)
	defer cancel()

	err := dynamodb.NewRepository().Health(ctx)
	if err != nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}

	err = message.New(ctx).Health()
	if err != nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}

	return c.NoContent(http.StatusOK)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}

	r, ok := i.(shared.EntityRequest)
	if !ok {
		return nil
	}

	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}
