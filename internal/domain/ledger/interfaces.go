package ledger

import "context"

//go:generate mockgen -source=interfaces.go -destination=mock.go -package=accounting
type Repository interface {
	SaveScript(ctx context.Context, s Config) error
	UpdateScript(ctx context.Context, s Config) error
	FindScriptByID(ctx context.Context, orgID string, scriptID string) (Config, error)
	FindScriptByLevel(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Config, error)
	FindAllScripts(ctx context.Context, orgID string, programID *int64) ([]Config, error)
}
