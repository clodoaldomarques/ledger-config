package ledger

import (
	"time"

	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
)

type ScriptResponse struct {
	ScriptID    string           `json:"script_id"`
	Level       string           `json:"level"`
	EventTypeID string           `json:"event_type_id"`
	OrgID       string           `json:"org_id"`
	ProgramID   *int64           `json:"program_id,omitempty"`
	Description string           `json:"description"`
	Company     *CompanyResponse `json:"company,omitempty"`
	Entries     []EntryResponse  `json:"entries"`
	Enable      bool             `json:"enable"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Version     int64            `json:"version"`
}

func buildScriptResponse(s ledger.Config) ScriptResponse {
	sr := ScriptResponse{
		ScriptID:    s.ConfigID,
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

	for _, e := range s.Entries {
		er := buildEntryResponse(e)
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

type EntryResponse struct {
	EntryTypeID   int64               `json:"entry_type_id"`
	Flow          string              `json:"flow"`
	Description   string              `json:"description"`
	AmountName    string              `json:"amount_name,omitempty"`
	Expression    string              `json:"expression,omitempty"`
	CashInBucket  string              `json:"cashin_bucket,omitempty"`
	CashOutBucket string              `json:"cashout_bucket,omitempty"`
	CostCenter    *CostCenterResponse `json:"cost_center,omitempty"`
	DebitAccount  *AccountResponse    `json:"debit_account,omitempty"`
	CreditAccount *AccountResponse    `json:"credit_account,omitempty"`
	Parameter     *ParameterResponse  `json:"parameter,omitempty"`
}

func buildEntryResponse(e ledger.Entry) EntryResponse {
	er := EntryResponse{
		EntryTypeID:   e.EntryTypeID,
		Flow:          e.Flow,
		Description:   e.Description,
		AmountName:    e.AmountName,
		Expression:    e.Expression,
		CashInBucket:  e.CashInBucket,
		CashOutBucket: e.CashOutBucket,
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

	if e.Parameter != nil {
		er.Parameter = buildParameterResponse(e.Parameter)
	}

	return er
}

type ParameterResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func buildParameterResponse(p *ledger.Parameter) *ParameterResponse {
	return &ParameterResponse{
		Name:  p.Name,
		Value: p.Value,
	}
}
