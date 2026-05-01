package server

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"gitlab.com/clodoaldomarques/accounting-scripts/internal/infra/rest/accounting"
	"gitlab.com/clodoaldomarques/accounting-scripts/internal/infra/rest/shared"
	"gitlab.com/clodoaldomarques/accounting-scripts/pkg/logger"
)

type Server struct {
	http *echo.Echo
}

func New() *Server {
	s := Server{
		http: echo.New(),
	}
	s.routes()
	return &s
}

func (s Server) routes() {
	s.http.Validator = &CustomValidator{validator: validator.New()}

	// health check
	s.http.GET("/", HealthCheck)

	// Accounting handler
	s.http.POST("/v1/accounting/scripts", accounting.CreateScript)
	s.http.PATCH("/v1/accounting/scripts/:script_id", accounting.UpdateScript)
	s.http.DELETE("/v1/accounting/scripts/:script_id/disable", accounting.DisableScript)

	s.http.GET("/v1/accounting/scripts/:event_type_id/:program_id", accounting.FindAccountingScript)
	s.http.GET("/v1/accounting/scripts", accounting.FindAllAccountingScripts)
}

func (s Server) Start(port int) error {
	return s.http.Start(fmt.Sprintf(":%d", port))
}

func HealthCheck(c echo.Context) error {
	logger.Info(c.Request().Context(), "health check", logger.Fields{})
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
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
