package ledger

import (
	"context"
	"fmt"
	"time"

	"github.com/clodoaldomarques/core-sdk/pkg/logger"
	"github.com/clodoaldomarques/core-sdk/pkg/tracer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

type Service struct {
	r Repository
}

func New(r Repository) *Service {
	return &Service{r: r}
}

func (s Service) CreateScript(ctx context.Context, cid string, scr Config) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::CreateScript", attribute.String("cid", cid))
	defer span.End()
	saved, err := s.r.FindConfigByLevel(ctx, string(scr.Level), scr.ProcessCode, scr.OrgID, &scr.ProgramID)
	if err == nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)
		return Config{}, fmt.Errorf("config was created with id: %v", saved.ConfigID)
	}

	if err := scr.Validate(); err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)

		logger.Error(ctx, "validate error", logger.Fields{
			"Error":  err.Error(),
			"Cid":    cid,
			"Config": scr,
		})
		return Config{}, err
	}

	scr.ConfigID = uuid.NewString()

	if err := s.r.SaveConfig(ctx, scr); err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)

		logger.Error(ctx, "error on save script", logger.Fields{
			"Error": err.Error(),
		})
		return Config{}, err
	}

	return scr, nil
}

func (s Service) UpdateScript(ctx context.Context, cid string, configID string, scr Config) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::UpdateScript", attribute.String("cid", cid))
	defer span.End()

	saved, err := s.r.FindConfigByID(ctx, scr.OrgID, configID)
	if err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)
		return Config{}, err
	}

	if !saved.Enable {
		return Config{}, ErrConfigNotFound{}
	}

	saved.Description = scr.Description
	saved.Scripts = scr.Scripts
	saved.UpdatedAt = time.Now().UTC()
	saved.Version++

	if err := saved.Validate(); err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)
		return Config{}, err
	}

	if err := s.r.UpdateConfig(ctx, saved); err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)
		return Config{}, err
	}

	return saved, nil
}

func (s Service) DisableScript(ctx context.Context, cid string, orgID string, scriptID string) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::DisableScript", attribute.String("cid", cid))
	defer span.End()
	saved, err := s.r.FindConfigByID(ctx, orgID, scriptID)
	if err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	if !saved.Enable {
		return Config{}, ErrConfigNotFound{}
	}

	saved.Enable = false
	saved.UpdatedAt = time.Now().UTC()
	saved.Version++

	if err := saved.Validate(); err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	if err := s.r.UpdateConfig(ctx, saved); err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	return saved, nil
}

func (s Service) EnableScript(ctx context.Context, cid string, orgID string, scriptID string) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::EnableScript", attribute.String("cid", cid))
	defer span.End()
	saved, err := s.r.FindConfigByID(ctx, orgID, scriptID)
	if err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	saved.Enable = true
	saved.UpdatedAt = time.Now().UTC()
	saved.Version++

	if err := saved.Validate(); err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	if err := s.r.UpdateConfig(ctx, saved); err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	return saved, nil
}

func (s Service) ActivateOrgID(ctx context.Context, cid string, orgID string) ([]Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::ActivateOrgID", attribute.String("cid", cid))
	defer span.End()
	configs, err := s.r.FindAllConfigs(ctx, "LEDGER", nil)
	if err != nil {
		span.AddAttributes(tracer.Attributes{
			"org_id": orgID,
		})
		span.SetError(err)
		return nil, err
	}

	newConfigs := make([]Config, 0, len(configs))

	for _, c := range configs {
		n := Config{
			ConfigID:    uuid.NewString(),
			Level:       TenantLevel,
			ProcessCode: c.ProcessCode,
			OrgID:       orgID,
			ProgramID:   c.ProgramID,
			Description: c.Description,
			Scripts:     c.Scripts,
			Enable:      true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Version:     1,
		}

		saved, _ := s.r.FindConfigByLevel(ctx, string(n.Level), n.ProcessCode, n.OrgID, &n.ProgramID)
		if saved.ConfigID == "" {
			err = s.r.SaveConfig(ctx, n)
			if err != nil {
				span.AddAttributes(tracer.Attributes{
					"org_id": orgID,
				})
				span.SetError(err)
				return nil, err
			}
			newConfigs = append(newConfigs, n)
		}
	}

	if len(newConfigs) == 0 {
		return nil, ErrOrgActivated{OrgID: orgID}
	}

	return newConfigs, nil
}

func (s Service) FindScriptByLevel(ctx context.Context, cid string, eventTypeID string, orgID string, programID int64) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::FindScriptByLevel", attribute.String("cid", cid))
	defer span.End()
	if saved, err := s.r.FindConfigByLevel(ctx, string(ProgramLevel), eventTypeID, orgID, &programID); err == nil && saved.Enable {
		span.AddAttributes(tracer.Attributes{
			"org_id":     orgID,
			"program_id": programID,
			"config":     saved,
		})
		span.SetError(err)
		logger.Info(ctx, "ledger config found",
			logger.Fields{
				"level":  string(ProgramLevel),
				"script": saved,
			})
		return saved, nil
	}

	if saved, err := s.r.FindConfigByLevel(ctx, string(TenantLevel), eventTypeID, orgID, &programID); err == nil && saved.Enable {
		span.AddAttributes(tracer.Attributes{
			"org_id":     orgID,
			"program_id": programID,
			"config":     saved,
		})
		span.SetError(err)
		logger.Info(ctx, "ledger config found",
			logger.Fields{
				"level":  string(ProgramLevel),
				"script": saved,
			})
		return saved, nil
	}

	return Config{}, ErrConfigNotFound{}
}

func (s Service) FindAllScripts(ctx context.Context, cid, orgID string, programID *int64) ([]Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::FindAllScripts", attribute.String("cid", cid))
	defer span.End()
	span.AddAttributes(tracer.Attributes{
		"org_id":     orgID,
		"program_id": programID,
	})
	return s.r.FindAllConfigs(ctx, orgID, programID)
}
