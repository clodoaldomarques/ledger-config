package ledger

import (
	"context"
	"fmt"
	"time"

	"github.com/clodoaldomarques/core-sdk/pkg/otel/tracer"
	"github.com/clodoaldomarques/core-sdk/pkg/zap/logger"
	"github.com/google/uuid"
)

type Service struct {
	r Repository
	t Topic
}

func New(r Repository, t Topic) *Service {
	return &Service{
		r: r,
		t: t,
	}
}

func (s Service) CreateConfig(ctx context.Context, cid string, scr Config) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::CreateConfig", map[string]any{
		"cid":    cid,
		"config": scr,
	})
	defer span.End()
	saved, err := s.r.FindConfigByLevel(ctx, cid, string(scr.Level), scr.ProcessingCode, scr.OrgID, &scr.ProgramID)
	if err == nil {
		span.AddAttributes(map[string]any{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(fmt.Errorf("config was created with id: %v", saved.ConfigID))
		return Config{}, fmt.Errorf("config was created with id: %v", saved.ConfigID)
	}

	if err := scr.Validate(); err != nil {
		span.AddAttributes(map[string]any{
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

	if err := s.r.SaveConfig(ctx, cid, scr); err != nil {
		span.AddAttributes(map[string]any{
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

	if err := s.t.Emit(ctx, cid, scr); err != nil {
		span.AddAttributes(map[string]any{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)
		return Config{}, err
	}

	return scr, nil
}

func (s Service) UpdateConfig(ctx context.Context, cid string, configID string, scr Config) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::UpdateScript", map[string]any{
		"cid":    cid,
		"config": scr,
	})
	defer span.End()

	saved, err := s.r.FindConfigByID(ctx, cid, scr.OrgID, configID)
	if err != nil {
		span.AddAttributes(map[string]any{
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
		span.AddAttributes(map[string]any{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)
		return Config{}, err
	}

	if err := s.r.UpdateConfig(ctx, cid, saved); err != nil {
		span.AddAttributes(map[string]any{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)
		return Config{}, err
	}

	if err := s.t.Emit(ctx, cid, scr); err != nil {
		span.AddAttributes(map[string]any{
			"org_id":     scr.OrgID,
			"program_id": scr.ProgramID,
			"config":     scr,
		})
		span.SetError(err)
		return Config{}, err
	}

	return saved, nil
}

func (s Service) DisableConfig(ctx context.Context, cid string, orgID string, scriptID string) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::DisableScript", map[string]any{
		"cid": cid,
	})
	defer span.End()
	saved, err := s.r.FindConfigByID(ctx, cid, orgID, scriptID)
	if err != nil {
		span.AddAttributes(map[string]any{
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
		span.AddAttributes(map[string]any{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	if err := s.r.UpdateConfig(ctx, cid, saved); err != nil {
		span.AddAttributes(map[string]any{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	if err := s.t.Emit(ctx, cid, saved); err != nil {
		span.AddAttributes(map[string]any{
			"org_id":     saved.OrgID,
			"program_id": saved.ProgramID,
			"config":     saved,
		})
		span.SetError(err)
		return Config{}, err
	}

	return saved, nil
}

func (s Service) EnableConfig(ctx context.Context, cid string, orgID string, scriptID string) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::EnableScript", map[string]any{
		"cid": cid,
	})
	defer span.End()
	saved, err := s.r.FindConfigByID(ctx, cid, orgID, scriptID)
	if err != nil {
		span.AddAttributes(map[string]any{
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
		span.AddAttributes(map[string]any{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	if err := s.r.UpdateConfig(ctx, cid, saved); err != nil {
		span.AddAttributes(map[string]any{
			"org_id":    orgID,
			"script_id": scriptID,
		})
		span.SetError(err)
		return Config{}, err
	}

	if err := s.t.Emit(ctx, cid, saved); err != nil {
		span.AddAttributes(map[string]any{
			"org_id":     saved.OrgID,
			"program_id": saved.ProgramID,
			"config":     saved,
		})
		span.SetError(err)
		return Config{}, err
	}

	return saved, nil
}

func (s Service) ActivateTenant(ctx context.Context, cid string, orgID string) ([]Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::ActivateTenant", map[string]any{
		"cid": cid,
	})
	defer span.End()
	configs, err := s.r.FindAllConfigs(ctx, cid, "LEDGER", nil)
	if err != nil {
		span.AddAttributes(map[string]any{
			"org_id": orgID,
		})
		span.SetError(err)
		return nil, err
	}

	newConfigs := make([]Config, 0, len(configs))

	for _, c := range configs {
		n := Config{
			ConfigID:       uuid.NewString(),
			Level:          TenantLevel,
			ProcessingCode: c.ProcessingCode,
			OrgID:          orgID,
			ProgramID:      c.ProgramID,
			Description:    c.Description,
			Scripts:        c.Scripts,
			Enable:         true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Version:        1,
		}

		saved, _ := s.r.FindConfigByLevel(ctx, cid, string(n.Level), n.ProcessingCode, n.OrgID, &n.ProgramID)
		if saved.ConfigID == "" {
			err = s.r.SaveConfig(ctx, cid, n)
			if err != nil {
				span.AddAttributes(map[string]any{
					"org_id": orgID,
				})
				span.SetError(err)
				return nil, err
			}

			if err := s.t.Emit(ctx, cid, n); err != nil {
				span.AddAttributes(map[string]any{
					"org_id":     n.OrgID,
					"program_id": n.ProgramID,
					"config":     n,
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

func (s Service) FindConfigByLevel(ctx context.Context, cid string, processingCode string, orgID string, programID int64) (Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::FindConfigByLevel", map[string]any{
		"cid":             cid,
		"processing_code": processingCode,
		"org_id":          orgID,
		"program_id":      programID,
	})
	defer span.End()

	if saved, err := s.r.FindConfigByLevel(ctx, cid, string(ProgramLevel), processingCode, orgID, &programID); err == nil && saved.Enable {
		span.AddAttributes(map[string]any{
			"org_id":     orgID,
			"program_id": programID,
			"config":     saved,
		})
		logger.Info(ctx, "ledger config found",
			logger.Fields{
				"level":  string(ProgramLevel),
				"script": saved,
			})
		return saved, nil
	}

	if saved, err := s.r.FindConfigByLevel(ctx, cid, string(TenantLevel), processingCode, orgID, &programID); err == nil && saved.Enable {
		span.AddAttributes(map[string]any{
			"org_id":     orgID,
			"program_id": programID,
			"config":     saved,
		})
		logger.Info(ctx, "ledger config found",
			logger.Fields{
				"level":  string(ProgramLevel),
				"script": saved,
			})
		return saved, nil
	}

	return Config{}, ErrConfigNotFound{}
}

func (s Service) FindAllConfigs(ctx context.Context, cid, orgID string, programID *int64) ([]Config, error) {
	span, ctx := tracer.NewSpanFromContext(ctx, "Service::FindAllConfigs", map[string]any{
		"cid": cid,
	})
	defer span.End()
	span.AddAttributes(map[string]any{
		"org_id":     orgID,
		"program_id": programID,
	})
	return s.r.FindAllConfigs(ctx, cid, orgID, programID)
}
