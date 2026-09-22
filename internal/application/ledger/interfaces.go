package ledger

import (
	"context"

	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
	"github.com/clodoaldomarques/ledger-config/internal/domain/program"
)

//go:generate mockgen -source=interfaces.go -destination=mock.go -package=ledger
type Repository interface {
	SaveConfig(ctx context.Context, cid string, s ledger.Config) error
	UpdateConfig(ctx context.Context, cid string, s ledger.Config) error
	FindConfigByID(ctx context.Context, cid string, orgID string, configID string) (ledger.Config, error)
	FindConfigByLevel(ctx context.Context, cid string, level string, eventTypeID string, orgID string, programID *int64) (ledger.Config, error)
	FindAllConfigs(ctx context.Context, cid string, orgID string, programID *int64) ([]ledger.Config, error)
}

type ProgramAPI interface {
	FindProgramaByID(ctx context.Context, orgID string, programID int) (program.Program, error)
	FindAllProgramaByOrgID(ctx context.Context, orgID string) ([]program.Program, error)
}

type Topic interface {
	Emit(ctx context.Context, cid string, e ledger.Config) error
}
