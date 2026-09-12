package ledger

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/clodoaldomarques/core-sdk/pkg/tracer"
	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
	"github.com/clodoaldomarques/ledger-config/internal/infra/db/dynamodb"
	"github.com/clodoaldomarques/ledger-config/internal/infra/message"
	"github.com/clodoaldomarques/ledger-config/internal/infra/rest/shared"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateConfig(c echo.Context) error {
	span, ctx, cid, orgID := NewSpanFromContext(c, "Handler::CreateConfig")
	defer span.End()

	r := dynamodb.NewRepository()
	defer r.Close()

	t := message.New(ctx)
	defer t.Close(ctx)

	s := ledger.New(r, t)

	psr := new(PostConfigRequest)
	if err := c.Bind(psr); err != nil {
		span.SetError(err)
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	if err := psr.Validate(); err != nil {
		span.SetError(err)
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	saved, err := s.CreateConfig(ctx, cid, psr.PostToEntity(orgID))
	if err != nil {
		span.SetError(err)
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	resp := buildConfigResponse(saved)

	return c.JSON(http.StatusCreated, resp)
}

func UpdateConfig(c echo.Context) error {
	span, ctx, cid, orgID := NewSpanFromContext(c, "Handler::UpdateConfig")
	defer span.End()

	configID := c.Param("config_id")

	r := dynamodb.NewRepository()
	defer r.Close()

	t := message.New(ctx)
	defer t.Close(ctx)

	s := ledger.New(r, t)

	psr := new(PathScriptRequest)
	if err := c.Bind(psr); err != nil {
		span.SetError(err)
		return echo.ErrBadRequest
	}

	if err := psr.Validate(); err != nil {
		span.SetError(err)
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	saved, err := s.UpdateConfig(ctx, cid, configID, psr.PatchToEntity(orgID))
	if err != nil {
		span.SetError(err)
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	resp := buildConfigResponse(saved)

	return c.JSON(http.StatusOK, resp)
}

func DisableConfig(c echo.Context) error {
	span, ctx, cid, orgID := NewSpanFromContext(c, "Handler::DisableConfig")
	defer span.End()

	scriptID := c.Param("script_id")
	r := dynamodb.NewRepository()
	defer r.Close()

	t := message.New(ctx)
	defer t.Close(ctx)

	s := ledger.New(r, t)

	saved, err := s.DisableConfig(ctx, cid, orgID, scriptID)
	if err != nil {
		span.SetError(err)
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}
	resp := buildConfigResponse(saved)

	return c.JSON(http.StatusOK, resp)
}

func ActivateTenant(c echo.Context) error {
	span, ctx, cid, _ := NewSpanFromContext(c, "Handler::ActivateTenant")
	defer span.End()

	r := dynamodb.NewRepository()
	defer r.Close()

	t := message.New(ctx)
	defer t.Close(ctx)

	s := ledger.New(r, t)

	poa := new(PostOrgActivate)
	if err := c.Bind(poa); err != nil {
		span.SetError(err)
		return echo.ErrBadRequest
	}

	if err := poa.Validate(); err != nil {
		span.SetError(err)
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	saved, err := s.ActivateTenant(ctx, cid, poa.OrgID)
	if err != nil {
		span.SetError(err)
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	resp := make([]ConfigResponse, 0, len(saved))
	for _, s := range saved {
		sr := buildConfigResponse(s)
		resp = append(resp, sr)
	}

	return c.JSON(http.StatusOK, resp)
}

func FindLedgerConfig(c echo.Context) error {
	span, ctx, cid, orgID := NewSpanFromContext(c, "Handler::FindLedgerConfig")
	defer span.End()

	evtID := strings.ToUpper(c.Param("event_type_id"))

	prgID, err := strconv.ParseInt(c.Param("program_id"), 10, 64)
	if err != nil {
		span.SetError(err)
		return c.JSON(http.StatusNotFound, shared.ErrResponse{Message: err.Error()})
	}

	r := dynamodb.NewRepository()
	defer r.Close()

	t := message.New(ctx)
	defer t.Close(ctx)

	s := ledger.New(r, t)

	scr, err := s.FindConfigByLevel(ctx, cid, evtID, orgID, prgID)
	if err != nil {
		span.SetError(err)
		return c.JSON(http.StatusNotFound, shared.ErrResponse{Message: err.Error()})
	}

	sr := buildConfigResponse(scr)

	return c.JSON(http.StatusOK, sr)
}

func FindAllLedgerConfig(c echo.Context) error {
	span, ctx, cid, orgID := NewSpanFromContext(c, "Handler::FindAllLedgerConfig")
	defer span.End()

	r := dynamodb.NewRepository()
	defer r.Close()

	t := message.New(ctx)
	defer t.Close(ctx)

	s := ledger.New(r, t)

	prgID := getProgramIDQueryParams(c)

	scrs, err := s.FindAllConfigs(ctx, cid, orgID, prgID)
	if err != nil {
		span.SetError(err)
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	resp := make([]ConfigResponse, 0, len(scrs))
	for _, s := range scrs {
		sr := buildConfigResponse(s)
		resp = append(resp, sr)
	}

	return c.JSON(http.StatusOK, resp)
}

func getProgramIDQueryParams(c echo.Context) *int64 {
	prg := c.QueryParam("program_id")
	prgID, err := strconv.ParseInt(prg, 10, 64)
	if err != nil {
		return nil
	}
	return &prgID
}

func NewSpanFromContext(e echo.Context, name string) (*tracer.TraceSpan, context.Context, string, string) {
	ctx := e.Request().Context()
	cid := e.Request().Header.Get("x-cid")
	orgID := e.Request().Header.Get("x-tenant")

	if cid == "" {
		cid = uuid.NewString()
	}

	span, ctx := tracer.NewSpanFromContext(ctx, name, map[string]any{
		"cid":    cid,
		"org_id": orgID,
	})

	return span, ctx, cid, orgID
}
