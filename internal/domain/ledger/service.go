package ledger

import (
	"context"
	"fmt"
	"time"

	"github.com/clodoaldomarques/ledger-config/pkg/logger"
	"github.com/google/uuid"
)

type Service struct {
	r Repository
}

func New(r Repository) *Service {
	return &Service{r: r}
}

func (s Service) CreateScript(ctx context.Context, cid string, scr Config) (Config, error) {
	saved, err := s.r.FindScriptByLevel(ctx, string(scr.Level), scr.EventTypeID, scr.OrgID, &scr.ProgramID)
	if err == nil {
		return Config{}, fmt.Errorf("script was created with id: %v", saved.ConfigID)
	}

	if err := scr.Validate(); err != nil {
		logger.Error(ctx, "validate error", logger.Fields{
			"Error":  err.Error(),
			"Cid":    cid,
			"Script": scr,
		})
		return Config{}, err
	}

	scr.ConfigID = uuid.NewString()

	if err := s.r.SaveScript(ctx, scr); err != nil {
		logger.Error(ctx, "error on save script", logger.Fields{
			"Error": err.Error(),
		})
		return Config{}, err
	}

	return scr, nil
}

func (s Service) UpdateScript(ctx context.Context, cid string, scriptID string, scr Config) (Config, error) {
	saved, err := s.r.FindScriptByID(ctx, scr.OrgID, scriptID)
	if err != nil {
		return Config{}, err
	}

	if !saved.Enable {
		return Config{}, ErrScriptNotFound{}
	}

	saved.Description = scr.Description
	saved.Company = scr.Company
	saved.Scripts = scr.Scripts
	saved.UpdatedAt = time.Now().UTC()
	saved.Version++

	if err := saved.Validate(); err != nil {
		return Config{}, err
	}

	if err := s.r.UpdateScript(ctx, saved); err != nil {
		return Config{}, err
	}

	return saved, nil
}

func (s Service) DisableScript(ctx context.Context, cid string, orgID string, scriptID string) (Config, error) {
	saved, err := s.r.FindScriptByID(ctx, orgID, scriptID)
	if err != nil {
		return Config{}, err
	}

	if !saved.Enable {
		return Config{}, ErrScriptNotFound{}
	}

	saved.Enable = false
	saved.UpdatedAt = time.Now().UTC()
	saved.Version++

	if err := saved.Validate(); err != nil {
		return Config{}, err
	}

	if err := s.r.UpdateScript(ctx, saved); err != nil {
		return Config{}, err
	}

	return saved, nil
}

func (s Service) EnableScript(ctx context.Context, cid string, orgID string, scriptID string) (Config, error) {
	saved, err := s.r.FindScriptByID(ctx, orgID, scriptID)
	if err != nil {
		return Config{}, err
	}

	saved.Enable = true
	saved.UpdatedAt = time.Now().UTC()
	saved.Version++

	if err := saved.Validate(); err != nil {
		return Config{}, err
	}

	if err := s.r.UpdateScript(ctx, saved); err != nil {
		return Config{}, err
	}

	return saved, nil
}

func (s Service) FindScriptByLevel(ctx context.Context, cid string, eventTypeID string, orgID string, programID int64) (Config, error) {
	if saved, err := s.r.FindScriptByLevel(ctx, string(ProgramLevel), eventTypeID, orgID, &programID); err == nil && saved.Enable {
		logger.Info(ctx, "accounting script found",
			logger.Fields{
				"level":  string(ProgramLevel),
				"script": saved,
			})
		return saved, nil
	}

	if saved, err := s.r.FindScriptByLevel(ctx, string(TenantLevel), eventTypeID, orgID, &programID); err == nil && saved.Enable {
		logger.Info(ctx, "accounting script found",
			logger.Fields{
				"level":  string(ProgramLevel),
				"script": saved,
			})
		return saved, nil
	}

	if saved, err := s.r.FindScriptByLevel(ctx, string(PlatformLevel), eventTypeID, "LEDGER", &programID); err == nil && saved.Enable {
		logger.Info(ctx, "accounting script found",
			logger.Fields{
				"level":  string(ProgramLevel),
				"script": saved,
			})
		return saved, nil
	}

	return Config{}, ErrScriptNotFound{}
}

func (s Service) FindAllScripts(ctx context.Context, cid, orgID string, programID *int64) ([]Config, error) {
	return s.r.FindAllScripts(ctx, orgID, programID)
}
