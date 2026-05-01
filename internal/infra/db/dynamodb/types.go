package dynamodb

import (
	"fmt"
	"time"

	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
)

type Config struct {
	OrgID       string    `dynamodbav:"org_id"`
	ConfigID    string    `dynamodbav:"config_id"`
	Filters     string    `dynamodbav:"filters"`
	Level       string    `dynamodbav:"level"`
	EventTypeID string    `dynamodbav:"event_type_id"`
	ProgramID   *int64    `dynamodbav:"program_id,omitempty"`
	Description string    `dynamodbav:"description_id"`
	Company     *Company  `dynamodbav:"company,omitempty"`
	Entries     []Script  `dynamodbav:"entries"`
	Enable      bool      `dynamodbav:"enable"`
	CreatedAt   time.Time `dynamodbav:"created_at"`
	UpdatedAt   time.Time `dynamodbav:"updated_at"`
	Version     int64     `dynamodbav:"version"`
}

func (s Config) toEntity() ledger.Config {
	scr := ledger.Config{
		ConfigID:    s.ConfigID,
		Level:       ledger.Level(s.Level),
		EventTypeID: s.EventTypeID,
		OrgID:       s.OrgID,
		Description: s.Description,
		Scripts:     make([]ledger.Script, 0),
		Enable:      s.Enable,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Version:     s.Version,
	}

	if s.ProgramID != nil {
		scr.ProgramID = *s.ProgramID
	}

	if s.Company != nil {
		scr.Company = s.Company.toEntity()
	}

	for _, e := range s.Entries {
		scr.Scripts = append(scr.Scripts, e.toEntity())
	}

	return scr
}

func buildConfigTable(s ledger.Config) Config {
	st := Config{
		OrgID:       s.OrgID,
		ConfigID:    s.ConfigID,
		Filters:     buildFilters(string(s.Level), s.OrgID, s.EventTypeID, s.ProgramID),
		Level:       string(s.Level),
		EventTypeID: s.EventTypeID,
		Description: s.Description,
		Entries:     make([]Script, 0),
		Enable:      s.Enable,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Version:     s.Version,
	}

	if s.ProgramID != 0 {
		st.ProgramID = &s.ProgramID
	}

	if s.Company != nil {
		st.Company = buildCompanyTable(s.Company)
	}

	for _, e := range s.Scripts {
		et := buildScriptTable(e)
		st.Entries = append(st.Entries, et)
	}

	return st
}

func buildFilters(level string, orgID string, eventTypeID string, programID int64) string {
	switch level {
	case string(ledger.PlatformLevel):
		return fmt.Sprintf("TENANT#LEDGER#PROGRAMID#0000#EVENTTYPEID#%s", eventTypeID)
	case string(ledger.TenantLevel):
		return fmt.Sprintf("TENANT#%s#PROGRAMID#0000#EVENTTYPEID#%s", orgID, eventTypeID)
	case string(ledger.ProgramLevel):
		return fmt.Sprintf("TENANT#%s#PROGRAMID#%04d#EVENTTYPEID#%s", orgID, programID, eventTypeID)
	default:
		return fmt.Sprintf("TENANT#%s#PROGRAMID#0000#EVENTTYPEID#%s", orgID, eventTypeID)
	}
}

func buildAllQuery(level string, orgID string, programID *int64) string {
	switch level {
	case string(ledger.PlatformLevel):
		return "TENANT#LEDGER#PROGRAMID#0000#"
	case string(ledger.TenantLevel):
		return fmt.Sprintf("TENANT#%s#PROGRAMID#0000#", orgID)
	case string(ledger.ProgramLevel):
		return fmt.Sprintf("TENANT#%s#PROGRAMID#%04d#", orgID, *programID)
	default:
		return fmt.Sprintf("TENANT#%s##PROGRAMID#0000#", orgID)
	}
}

type Company struct {
	Code string `dynamodbav:"code"`
	Type string `dynamodbav:"type"`
}

func buildCompanyTable(c *ledger.Company) *Company {
	return &Company{
		Code: c.Code,
		Type: c.Type,
	}
}

func (c Company) toEntity() *ledger.Company {
	return &ledger.Company{
		Code: c.Code,
		Type: c.Type,
	}
}

type CostCenter struct {
	DebitCost  string `dynamodbav:"debit_cost"`
	DebitOrg   string `dynamodbav:"debit_org"`
	CreditCost string `dynamodbav:"credit_cost"`
	CreditOrg  string `dynamodbav:"credit_org"`
}

func buildCostCenterTable(c *ledger.CostCenter) *CostCenter {
	return &CostCenter{
		DebitCost:  c.DebitCost,
		DebitOrg:   c.DebitOrg,
		CreditCost: c.CreditCost,
		CreditOrg:  c.CreditOrg,
	}
}

func (c CostCenter) toEntity() *ledger.CostCenter {
	return &ledger.CostCenter{
		DebitCost:  c.DebitCost,
		DebitOrg:   c.DebitOrg,
		CreditCost: c.CreditCost,
		CreditOrg:  c.CreditOrg,
	}
}

type Account struct {
	Type        string `dynamodbav:"type"`
	Number      string `dynamodbav:"number"`
	Description string `dynamodbav:"description"`
	Cosif       string `dynamodbav:"cosif"`
}

func buildAccountTable(a *ledger.Account) *Account {
	return &Account{
		Number:      a.Number,
		Description: a.Description,
		Cosif:       a.Cosif,
	}
}

func (a Account) toEntity() *ledger.Account {
	return &ledger.Account{
		Number:      a.Number,
		Description: a.Description,
		Cosif:       a.Cosif,
	}
}

type Script struct {
	ScriptID      int64       `dynamodbav:"script_id"`
	Flow          string      `dynamodbav:"flow"`
	Description   string      `dynamodbav:"description"`
	AmountName    string      `dynamodbav:"amount_name"`
	Expression    string      `dynamodbav:"expression"`
	CostCenter    *CostCenter `dynamodbav:"cost_center,omitempty"`
	DebitAccount  *Account    `dynamodbav:"debit_account,omitempty"`
	CreditAccount *Account    `dynamodbav:"credit_account,omitempty"`
}

func buildScriptTable(e ledger.Script) Script {
	et := Script{
		ScriptID:    e.ScriptID,
		Flow:        e.Flow,
		Description: e.Description,
		Expression:  e.Expression,
	}

	if e.CostCenter != nil {
		et.CostCenter = buildCostCenterTable(e.CostCenter)
	}

	if e.DebitAccount != nil {
		et.DebitAccount = buildAccountTable(e.DebitAccount)
	}

	if e.CreditAccount != nil {
		et.CreditAccount = buildAccountTable(e.CreditAccount)
	}

	return et
}

func (e Script) toEntity() ledger.Script {
	et := ledger.Script{
		ScriptID:    e.ScriptID,
		Flow:        e.Flow,
		Description: e.Description,
		Expression:  e.Expression,
	}

	if e.CostCenter != nil {
		et.CostCenter = e.CostCenter.toEntity()
	}

	if e.CreditAccount != nil {
		et.CreditAccount = e.CreditAccount.toEntity()
	}

	if e.DebitAccount != nil {
		et.DebitAccount = e.DebitAccount.toEntity()
	}

	return et
}
