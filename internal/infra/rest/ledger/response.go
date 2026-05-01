package ledger

import (
	"time"

	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
)

type ConfigResponse struct {
	ConfigID    string           `json:"config_id"`
	Level       string           `json:"level"`
	EventTypeID string           `json:"event_type_id"`
	OrgID       string           `json:"org_id"`
	ProgramID   *int64           `json:"program_id,omitempty"`
	Description string           `json:"description"`
	Company     *CompanyResponse `json:"company,omitempty"`
	Entries     []ScriptResponse `json:"entries"`
	Enable      bool             `json:"enable"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Version     int64            `json:"version"`
}

func buildConfigResponse(s ledger.Config) ConfigResponse {
	sr := ConfigResponse{
		ConfigID:    s.ConfigID,
		Level:       string(s.Level),
		EventTypeID: s.EventTypeID,
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

	if s.Company != nil {
		sr.Company = buildCompanyResponse(s.Company)
	}

	for _, e := range s.Scripts {
		er := buildScriptResponse(e)
		sr.Entries = append(sr.Entries, er)
	}

	return sr

}

type CompanyResponse struct {
	Code string `json:"code,omitempty"`
	Type string `json:"type,omitempty"`
}

func buildCompanyResponse(c *ledger.Company) *CompanyResponse {
	return &CompanyResponse{
		Code: c.Code,
		Type: c.Type,
	}
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

type CostCenterResponse struct {
	DebitCost  string `json:"debit_cost"`
	DebitOrg   string `json:"debit_org"`
	CreditCost string `json:"credit_cost"`
	CreditOrg  string `json:"credit_org"`
}

func buildCostCenterResponse(c *ledger.CostCenter) *CostCenterResponse {
	return &CostCenterResponse{
		DebitCost:  c.DebitCost,
		DebitOrg:   c.DebitOrg,
		CreditCost: c.CreditCost,
		CreditOrg:  c.CreditOrg,
	}

}

type ScriptResponse struct {
	ScriptID      int64               `json:"script_id"`
	Flow          string              `json:"flow"`
	Description   string              `json:"description"`
	Expression    string              `json:"expression,omitempty"`
	CostCenter    *CostCenterResponse `json:"cost_center,omitempty"`
	DebitAccount  *AccountResponse    `json:"debit_account,omitempty"`
	CreditAccount *AccountResponse    `json:"credit_account,omitempty"`
}

func buildScriptResponse(e ledger.Script) ScriptResponse {
	er := ScriptResponse{
		ScriptID:    e.ScriptID,
		Flow:        e.Flow,
		Description: e.Description,
		Expression:  e.Expression,
	}

	if e.CostCenter != nil {
		er.CostCenter = buildCostCenterResponse(e.CostCenter)
	}

	if e.DebitAccount != nil {
		er.DebitAccount = buildAccountResponse(e.DebitAccount)
	}

	if e.CreditAccount != nil {
		er.CreditAccount = buildAccountResponse(e.CreditAccount)
	}

	return er
}
