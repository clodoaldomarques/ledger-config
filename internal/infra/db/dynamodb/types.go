package dynamodb

import (
	"fmt"
	"time"

	"gitlab.com/clodoaldomarques/accounting-scripts/internal/domain/accounting"
)

type Script struct {
	OrgID       string    `dynamodbav:"org_id"`
	ScriptID    string    `dynamodbav:"script_id"`
	Filters     string    `dynamodbav:"filters"`
	Level       string    `dynamodbav:"level"`
	EventTypeID string    `dynamodbav:"event_type_id"`
	ProgramID   *int64    `dynamodbav:"program_id,omitempty"`
	Description string    `dynamodbav:"description_id"`
	Company     *Company  `dynamodbav:"company,omitempty"`
	Entries     []Entry   `dynamodbav:"entries"`
	Enable      bool      `dynamodbav:"enable"`
	CreatedAt   time.Time `dynamodbav:"created_at"`
	UpdatedAt   time.Time `dynamodbav:"updated_at"`
	Version     int64     `dynamodbav:"version"`
}

func (s Script) toEntity() accounting.Script {
	scr := accounting.Script{
		ScriptID:    s.ScriptID,
		Level:       accounting.Level(s.Level),
		EventTypeID: s.EventTypeID,
		OrgID:       s.OrgID,
		Description: s.Description,
		Entries:     make([]accounting.Entry, 0),
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
		scr.Entries = append(scr.Entries, e.toEntity())
	}

	return scr
}

func buildScriptTable(s accounting.Script) Script {
	st := Script{
		OrgID:       s.OrgID,
		ScriptID:    s.ScriptID,
		Filters:     buildFilters(string(s.Level), s.OrgID, s.EventTypeID, s.ProgramID),
		Level:       string(s.Level),
		EventTypeID: s.EventTypeID,
		Description: s.Description,
		Entries:     make([]Entry, 0),
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

	for _, e := range s.Entries {
		et := buildEntryTable(e)
		st.Entries = append(st.Entries, et)
	}

	return st
}

func buildFilters(level string, orgID string, eventTypeID string, programID int64) string {
	switch level {
	case string(accounting.PlatformLevel):
		return fmt.Sprintf("TENANT#PISMO#PROGRAMID#0000#EVENTTYPEID#%s", eventTypeID)
	case string(accounting.TenantLevel):
		return fmt.Sprintf("TENANT#%s#PROGRAMID#0000#EVENTTYPEID#%s", orgID, eventTypeID)
	case string(accounting.ProgramLevel):
		return fmt.Sprintf("TENANT#%s#PROGRAMID#%04d#EVENTTYPEID#%s", orgID, programID, eventTypeID)
	default:
		return fmt.Sprintf("TENANT#%s#PROGRAMID#0000#EVENTTYPEID#%s", orgID, eventTypeID)
	}
}

func buildAllQuery(level string, orgID string, programID *int64) string {
	switch level {
	case string(accounting.PlatformLevel):
		return "TENANT#PISMO#PROGRAMID#0000#"
	case string(accounting.TenantLevel):
		return fmt.Sprintf("TENANT#%s#PROGRAMID#0000#", orgID)
	case string(accounting.ProgramLevel):
		return fmt.Sprintf("TENANT#%s#PROGRAMID#%04d#", orgID, *programID)
	default:
		return fmt.Sprintf("TENANT#%s##PROGRAMID#0000#", orgID)
	}
}

type Company struct {
	Code string `dynamodbav:"code"`
	Type string `dynamodbav:"type"`
}

func buildCompanyTable(c *accounting.Company) *Company {
	return &Company{
		Code: c.Code,
		Type: c.Type,
	}
}

func (c Company) toEntity() *accounting.Company {
	return &accounting.Company{
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

func buildCostCenterTable(c *accounting.CostCenter) *CostCenter {
	return &CostCenter{
		DebitCost:  c.DebitCost,
		DebitOrg:   c.DebitOrg,
		CreditCost: c.CreditCost,
		CreditOrg:  c.CreditOrg,
	}
}

func (c CostCenter) toEntity() *accounting.CostCenter {
	return &accounting.CostCenter{
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

func buildAccountTable(a *accounting.Account) *Account {
	return &Account{
		Number:      a.Number,
		Description: a.Description,
		Cosif:       a.Cosif,
	}
}

func (a Account) toEntity() *accounting.Account {
	return &accounting.Account{
		Number:      a.Number,
		Description: a.Description,
		Cosif:       a.Cosif,
	}
}

type Entry struct {
	EntryTypeID   int64       `dynamodbav:"entry_type_id"`
	Flow          string      `dynamodbav:"flow"`
	Description   string      `dynamodbav:"description"`
	AmountName    string      `dynamodbav:"amount_name"`
	Expression    string      `dynamodbav:"expression"`
	CashInBucket  string      `dynamodbav:"cashin_bucket"`
	CashOutBucket string      `dynamodbav:"cashout_bucket"`
	CostCenter    *CostCenter `dynamodbav:"cost_center,omitempty"`
	DebitAccount  *Account    `dynamodbav:"debit_account,omitempty"`
	CreditAccount *Account    `dynamodbav:"credit_account,omitempty"`
	Parameter     *Parameter  `dynamodbav:"parameter,omitempty"`
}

func buildEntryTable(e accounting.Entry) Entry {
	et := Entry{
		EntryTypeID:   e.EntryTypeID,
		Flow:          e.Flow,
		Description:   e.Description,
		AmountName:    e.AmountName,
		Expression:    e.Expression,
		CashInBucket:  e.CashInBucket,
		CashOutBucket: e.CashOutBucket,
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

	if e.Parameter != nil {
		et.Parameter = buildParameterTable(e.Parameter)
	}

	return et
}

func (e Entry) toEntity() accounting.Entry {
	et := accounting.Entry{
		EntryTypeID:   e.EntryTypeID,
		Flow:          e.Flow,
		Description:   e.Description,
		AmountName:    e.AmountName,
		Expression:    e.Expression,
		CashInBucket:  e.CashInBucket,
		CashOutBucket: e.CashOutBucket,
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

	if e.Parameter != nil {
		et.Parameter = e.Parameter.toEntity()
	}

	return et
}

type Parameter struct {
	Name  string `dynamodbav:"name"`
	Value string `dynamodbav:"value"`
}

func buildParameterTable(p *accounting.Parameter) *Parameter {
	return &Parameter{
		Name:  p.Name,
		Value: p.Value,
	}
}

func (p Parameter) toEntity() *accounting.Parameter {
	return &accounting.Parameter{
		Name:  p.Name,
		Value: p.Value,
	}
}
