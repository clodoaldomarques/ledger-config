package ledger

import (
	"fmt"
	"time"
)

type Level string

const (
	PlatformLevel Level = "platform"
	TenantLevel   Level = "tenant"
	ProgramLevel  Level = "program"
)

const (
	PlatformTenant = "LEDGER"
)

type Config struct {
	ConfigID    string
	Level       Level
	ProcessCode string
	OrgID       string
	ProgramID   int64
	Description string
	Scripts     []Script
	Enable      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int64
}

func (c Config) Validate() error {
	switch c.Level {
	case TenantLevel:
		if c.ProgramID != 0 {
			return fmt.Errorf("program_id is not required for level: %s", string(TenantLevel))
		}
		if c.OrgID == "" {
			return fmt.Errorf("org_id is required for level: %s", string(TenantLevel))
		}
	case ProgramLevel:
		if c.OrgID == "" || c.ProgramID == 0 {
			return fmt.Errorf("org_id and program_id are required for level: %s", string(ProgramLevel))
		}
	case PlatformLevel:
		if c.OrgID != PlatformTenant {
			return fmt.Errorf("this tenant %s can not create %s level script", c.OrgID, c.Level)
		}
	default:
		return fmt.Errorf("invalid script level %v, chose tenant or program", c.Level)
	}

	if err := validateScript(c.Scripts); err != nil {
		return err
	}
	return nil
}

func validateScript(scripts []Script) error {
	m := make(map[string]any)
	for _, s := range scripts {
		if _, has := m[s.ScriptKey()]; has {
			return ErrDuplicatedScript{msg: fmt.Sprintf("duplicated script: %1s - %2d - %3s", s.Flow, s.ScriptID, s.Description)}
		}
		m[s.ScriptKey()] = s

		if err := s.Validate(); err != nil {
			return err
		}
	}

	return nil
}
