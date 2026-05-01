package accounting

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/clodoaldomarques/accounting-scripts/pkg/logger"
)

type Service struct {
	r Repository
}

func New(r Repository) *Service {
	return &Service{r: r}
}

func (s Service) CreateScript(ctx context.Context, cid string, scr Script) (Script, error) {
	saved, err := s.r.FindScriptByLevel(ctx, string(scr.Level), scr.EventTypeID, scr.OrgID, &scr.ProgramID)
	if err == nil {
		return Script{}, fmt.Errorf("script was created with id: %v", saved.ScriptID)
	}

	if err := scr.Validate(); err != nil {
		logger.Error(ctx, "validate error", logger.Fields{
			"Error":  err.Error(),
			"Cid":    cid,
			"Script": scr,
		})
		return Script{}, err
	}

	scr.ScriptID = uuid.NewString()

	if err := s.r.SaveScript(ctx, scr); err != nil {
		logger.Error(ctx, "error on save script", logger.Fields{
			"Error": err.Error(),
		})
		return Script{}, err
	}

	return scr, nil
}

func (s Service) UpdateScript(ctx context.Context, cid string, scriptID string, scr Script) (Script, error) {
	saved, err := s.r.FindScriptByID(ctx, scr.OrgID, scriptID)
	if err != nil {
		return Script{}, err
	}

	if !saved.Enable {
		return Script{}, ErrScriptNotFound{}
	}

	saved.Description = scr.Description
	saved.Company = scr.Company
	saved.Entries = scr.Entries
	saved.UpdatedAt = time.Now().UTC()
	saved.Version++

	if err := saved.Validate(); err != nil {
		return Script{}, err
	}

	if err := s.r.UpdateScript(ctx, saved); err != nil {
		return Script{}, err
	}

	return saved, nil
}

func (s Service) DisableScript(ctx context.Context, cid string, orgID string, scriptID string) (Script, error) {
	saved, err := s.r.FindScriptByID(ctx, orgID, scriptID)
	if err != nil {
		return Script{}, err
	}

	if !saved.Enable {
		return Script{}, ErrScriptNotFound{}
	}

	saved.Enable = false
	saved.UpdatedAt = time.Now().UTC()
	saved.Version++

	if err := saved.Validate(); err != nil {
		return Script{}, err
	}

	if err := s.r.UpdateScript(ctx, saved); err != nil {
		return Script{}, err
	}

	return saved, nil
}

func (s Service) EnableScript(ctx context.Context, cid string, orgID string, scriptID string) (Script, error) {
	saved, err := s.r.FindScriptByID(ctx, orgID, scriptID)
	if err != nil {
		return Script{}, err
	}

	saved.Enable = true
	saved.UpdatedAt = time.Now().UTC()
	saved.Version++

	if err := saved.Validate(); err != nil {
		return Script{}, err
	}

	if err := s.r.UpdateScript(ctx, saved); err != nil {
		return Script{}, err
	}

	return saved, nil
}

func (s Service) FindScriptByLevel(ctx context.Context, cid string, eventTypeID string, orgID string, programID int64) (Script, error) {
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

	if saved, err := s.r.FindScriptByLevel(ctx, string(PlatformLevel), eventTypeID, "PISMO", &programID); err == nil && saved.Enable {
		logger.Info(ctx, "accounting script found",
			logger.Fields{
				"level":  string(ProgramLevel),
				"script": saved,
			})
		return saved, nil
	}

	return Script{}, ErrScriptNotFound{}
}

func (s Service) FindAllScripts(ctx context.Context, cid, orgID string, programID *int64) ([]Script, error) {
	return s.r.FindAllScripts(ctx, orgID, programID)
}
