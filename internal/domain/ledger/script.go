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
	EventTypeID string
	OrgID       string
	ProgramID   int64
	Description string
	Company     *Company
	Entries     []Entry
	Enable      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int64
}

func (s Config) Validate() error {
	switch s.Level {
	case TenantLevel:
		if s.ProgramID != 0 {
			return fmt.Errorf("program_id is not required for level: %s", string(TenantLevel))
		}
		if s.OrgID == "" {
			return fmt.Errorf("org_id is required for level: %s", string(TenantLevel))
		}
	case ProgramLevel:
		if s.OrgID == "" || s.ProgramID == 0 {
			return fmt.Errorf("org_id and program_id are required for level: %s", string(ProgramLevel))
		}
	case PlatformLevel:
		if s.OrgID != PlatformTenant {
			return fmt.Errorf("this tenant %s can not create %s level script", s.OrgID, s.Level)
		}
		if s.Company != nil {
			return fmt.Errorf("company not required to %v level", PlatformLevel)
		}
	default:
		return fmt.Errorf("invalid script level %v, chose tenant or program", s.Level)
	}

	if err := validateEntry(s.Entries); err != nil {
		return err
	}
	return nil
}

func validateEntry(entries []Entry) error {
	m := make(map[int64]any)
	for _, e := range entries {
		if _, has := m[e.EntryTypeID]; has {
			return ErrDuplicatedEntry{msg: fmt.Sprintf("duplicated entry: %d - %s", e.EntryTypeID, e.Description)}
		}
		m[e.EntryTypeID] = e

		if err := e.Validate(); err != nil {
			return err
		}
	}

	return nil
}
