package ledger

import "context"

//go:generate mockgen -source=interfaces.go -destination=mock.go -package=ledger
type Repository interface {
	SaveConfig(ctx context.Context, s Config) error
	UpdateConfig(ctx context.Context, s Config) error
	FindConfigByID(ctx context.Context, orgID string, configID string) (Config, error)
	FindConfigByLevel(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Config, error)
	FindAllConfigs(ctx context.Context, orgID string, programID *int64) ([]Config, error)
}
