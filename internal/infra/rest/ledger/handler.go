package ledger

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
	"github.com/clodoaldomarques/ledger-config/internal/infra/db/dynamodb"
	"github.com/clodoaldomarques/ledger-config/internal/infra/rest/shared"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateConfig(c echo.Context) error {
	orgID, cid := getHeaders(c)
	ctx := c.Request().Context()

	r := dynamodb.NewRepository()
	defer r.Close()

	s := ledger.New(r)

	psr := new(PostConfigRequest)
	if err := c.Bind(psr); err != nil {
		return echo.ErrBadRequest
	}

	if err := psr.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	saved, err := s.CreateScript(ctx, cid, psr.PostToEntity(orgID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	resp := buildConfigResponse(saved)

	return c.JSON(http.StatusCreated, resp)
}

func UpdateConfig(c echo.Context) error {
	orgID, cid := getHeaders(c)
	configID := c.Param("config_id")
	ctx := c.Request().Context()
	r := dynamodb.NewRepository()
	defer r.Close()

	s := ledger.New(r)

	psr := new(PathScriptRequest)
	if err := c.Bind(psr); err != nil {
		return echo.ErrBadRequest
	}

	if err := psr.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	saved, err := s.UpdateScript(ctx, cid, configID, psr.PatchToEntity(orgID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	resp := buildConfigResponse(saved)

	return c.JSON(http.StatusOK, resp)
}

func DisableConfig(c echo.Context) error {
	orgID, cid := getHeaders(c)
	scriptID := c.Param("script_id")
	ctx := c.Request().Context()
	r := dynamodb.NewRepository()
	defer r.Close()

	s := ledger.New(r)

	saved, err := s.DisableScript(ctx, cid, orgID, scriptID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}
	resp := buildConfigResponse(saved)

	return c.JSON(http.StatusOK, resp)
}

func FindLedgerConfig(c echo.Context) error {
	ctx := c.Request().Context()
	orgID, cid := getHeaders(c)
	evtID := strings.ToUpper(c.Param("event_type_id"))

	prgID, err := strconv.ParseInt(c.Param("program_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusNotFound, shared.ErrResponse{Message: err.Error()})
	}

	r := dynamodb.NewRepository()
	defer r.Close()

	s := ledger.New(r)

	scr, err := s.FindScriptByLevel(ctx, cid, evtID, orgID, prgID)
	if err != nil {
		return c.JSON(http.StatusNotFound, shared.ErrResponse{Message: err.Error()})
	}

	sr := buildConfigResponse(scr)

	return c.JSON(http.StatusOK, sr)
}

func FindAllLedgerConfig(c echo.Context) error {
	orgID, cid := getHeaders(c)
	ctx := c.Request().Context()

	r := dynamodb.NewRepository()
	defer r.Close()

	s := ledger.New(r)

	prgID := getProgramIDQueryParams(c)

	scrs, err := s.FindAllScripts(ctx, cid, orgID, prgID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.ErrResponse{Message: err.Error()})
	}

	resp := make([]ConfigResponse, 0, len(scrs))
	for _, s := range scrs {
		sr := buildConfigResponse(s)
		resp = append(resp, sr)
	}

	return c.JSON(http.StatusOK, resp)
}

func getHeaders(c echo.Context) (string, string) {
	orgID := c.Request().Header.Get("x-tenant")
	cid := c.Request().Header.Get("x-cid")

	if cid == "" {
		cid = uuid.NewString()
	}
	return orgID, cid
}

func getProgramIDQueryParams(c echo.Context) *int64 {
	prg := c.QueryParam("program_id")
	prgID, err := strconv.ParseInt(prg, 10, 64)
	if err != nil {
		return nil
	}
	return &prgID
}
