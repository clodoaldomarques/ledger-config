package accounting

import "context"

//go:generate mockgen -source=interfaces.go -destination=mock.go -package=accounting
type Repository interface {
	SaveScript(ctx context.Context, s Script) error
	UpdateScript(ctx context.Context, s Script) error
	FindScriptByID(ctx context.Context, orgID string, scriptID string) (Script, error)
	FindScriptByLevel(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Script, error)
	FindAllScripts(ctx context.Context, orgID string, programID *int64) ([]Script, error)
}
