package ledger

import (
	"context"

	"github.com/clodoaldomarques/ledger-config/internal/domain/program"
)

//go:generate mockgen -source=interfaces.go -destination=mock.go -package=ledger
type Repository interface {
	SaveConfig(ctx context.Context, s Config) error
	UpdateConfig(ctx context.Context, s Config) error
	FindConfigByID(ctx context.Context, orgID string, configID string) (Config, error)
	FindConfigByLevel(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Config, error)
	FindAllConfigs(ctx context.Context, orgID string, programID *int64) ([]Config, error)
}

type ProgramAPI interface {
	FindProgramaByID(ctx context.Context, orgID string, programID int) (program.Program, error)
	FindAllProgramaByOrgID(ctx context.Context, orgID string) ([]program.Program, error)
}

type Topic interface {
	Emit(ctx context.Context, cid string, e Config) error
}
