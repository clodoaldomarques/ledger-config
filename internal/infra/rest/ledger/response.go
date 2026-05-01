package ledger

import (
	"time"

	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
)

type ConfigResponse struct {
	ConfigID    string           `json:"config_id"`
	Level       string           `json:"level"`
	ProcessCode string           `json:"process_code"`
	OrgID       string           `json:"org_id"`
	ProgramID   *int64           `json:"program_id,omitempty"`
	Description string           `json:"description"`
	Scripts     []ScriptResponse `json:"scripts"`
	Enable      bool             `json:"enable"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Version     int64            `json:"version"`
}

func buildConfigResponse(s ledger.Config) ConfigResponse {
	sr := ConfigResponse{
		ConfigID:    s.ConfigID,
		Level:       string(s.Level),
		ProcessCode: s.ProcessCode,
		OrgID:       s.OrgID,
		Description: s.Description,
		Enable:      s.Enable,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Version:     s.Version,
	}

	if s.ProgramID > 0 {
		sr.ProgramID = &s.ProgramID
	}

	for _, e := range s.Scripts {
		er := buildScriptResponse(e)
		sr.Scripts = append(sr.Scripts, er)
	}

	return sr

}

type AccountResponse struct {
	Number      string `json:"number"`
	Description string `json:"description"`
	Cosif       string `json:"cosif,omitempty"`
}

func buildAccountResponse(a *ledger.Account) *AccountResponse {
	return &AccountResponse{
		Number:      a.Number,
		Description: a.Description,
		Cosif:       a.Cosif,
	}
}

type ScriptResponse struct {
	ScriptID      int64            `json:"script_id"`
	Flow          string           `json:"flow"`
	Description   string           `json:"description"`
	Expression    string           `json:"expression,omitempty"`
	DebitAccount  *AccountResponse `json:"debit_account,omitempty"`
	CreditAccount *AccountResponse `json:"credit_account,omitempty"`
}

func buildScriptResponse(e ledger.Script) ScriptResponse {
	er := ScriptResponse{
		ScriptID:    e.ScriptID,
		Flow:        e.Flow,
		Description: e.Description,
		Expression:  e.Expression,
	}

	if e.DebitAccount != nil {
		er.DebitAccount = buildAccountResponse(e.DebitAccount)
	}

	if e.CreditAccount != nil {
		er.CreditAccount = buildAccountResponse(e.CreditAccount)
	}

	return er
}
