package accounting

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/clodoaldomarques/ledger-config/internal/domain/accounting"
)

var (
	validAmounts = []string{"amount", "interest", "duo_date"}
)

type PostScriptRequest struct {
	Level       string             `json:"level" validate:"required"`
	EventTypeID string             `json:"event_type_id" validate:"required"`
	ProgramID   *int64             `json:"program_id,omitempty"`
	Description string             `json:"description" validate:"required"`
	Company     *CompanyRequest    `json:"company,omitempty"`
	CostCenter  *CostCenterRequest `json:"cost_center,omitempty"`
	Accounts    []AccountRequest   `json:"accounts,omitempty"`
	Entries     []EntryRequest     `json:"entries" validate:"required"`
	Enable      *bool              `json:"enable,omitempty"`
}

func (p PostScriptRequest) Validate() error {
	if p.Level == "" {
		return fmt.Errorf("level is required")
	}

	if p.Description == "" {
		return fmt.Errorf("description is required")
	}

	if len(p.Entries) == 0 {
		return fmt.Errorf("entries is required")
	}

	for _, e := range p.Entries {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (p PostScriptRequest) PostToEntity(orgID string) accounting.Script {
	scr := accounting.Script{
		Level:       accounting.Level(p.Level),
		EventTypeID: strings.ToUpper(p.EventTypeID),
		OrgID:       orgID,
		Description: p.Description,
		Entries:     make([]accounting.Entry, 0, len(p.Entries)),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Enable:      true,
		Version:     1,
	}

	if p.ProgramID != nil {
		scr.ProgramID = *p.ProgramID
	}

	if p.Company != nil {
		scr.Company = p.Company.ToEntity()
	}

	for _, e := range p.Entries {
		scr.Entries = append(scr.Entries, e.ToEntity())
	}

	return scr
}

type PathScriptRequest struct {
	Description string          `json:"description" validate:"required"`
	Company     *CompanyRequest `json:"company,omitempty"`
	Entries     []EntryRequest  `json:"entries" validate:"required"`
	Enable      *bool           `json:"enable,omitempty"`
}

func (p PathScriptRequest) Validate() error {

	if p.Description == "" {
		return fmt.Errorf("description is required")
	}

	if len(p.Entries) == 0 {
		return fmt.Errorf("entries is required")
	}

	for _, e := range p.Entries {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (p PathScriptRequest) PatchToEntity(orgID string) accounting.Script {
	scr := accounting.Script{
		OrgID:       orgID,
		Description: p.Description,
		Entries:     make([]accounting.Entry, 0, len(p.Entries)),
	}

	if p.Company != nil {
		scr.Company = p.Company.ToEntity()
	}

	for _, e := range p.Entries {
		scr.Entries = append(scr.Entries, e.ToEntity())
	}

	if p.Enable != nil {
		scr.Enable = *p.Enable
	}

	return scr
}

type CompanyRequest struct {
	Code string `json:"code,omitempty"`
	Type string `json:"type,omitempty"`
}

func (c CompanyRequest) ToEntity() *accounting.Company {
	return &accounting.Company{
		Code: c.Code,
		Type: c.Type,
	}
}

type AccountRequest struct {
	Number      string `json:"number" validate:"required"`
	Description string `json:"description" validate:"required"`
	Cosif       string `json:"cosif,omitempty"`
}

func (a AccountRequest) Validate() error {
	if a.Number == "" {
		return errors.New("entry.account.number is required")
	}
	if a.Description == "" {
		return errors.New("entry.account.description is required")
	}
	return nil
}

func (a AccountRequest) ToEntity() *accounting.Account {
	return &accounting.Account{
		Number:      a.Number,
		Description: a.Description,
		Cosif:       a.Cosif,
	}
}

type CostCenterRequest struct {
	DebitCost  string `json:"debit_cost" validate:"required"`
	DebitOrg   string `json:"debit_org" validate:"required"`
	CreditCost string `json:"credit_cost" validate:"required"`
	CreditOrg  string `json:"credit_org" validate:"required"`
}

func (c CostCenterRequest) Validate() error {
	if c.DebitCost == "" {
		return errors.New("entry.cost_center.debit_cost is required")
	}
	if c.DebitOrg == "" {
		return errors.New("entry.cost_center.debit_org is required")
	}
	if c.CreditCost == "" {
		return errors.New("entry.cost_center.credit_cost is required")
	}
	if c.CreditOrg == "" {
		return errors.New("entry.cost_center.credit_org is required")
	}
	return nil
}

func (c CostCenterRequest) ToEntity() *accounting.CostCenter {
	return &accounting.CostCenter{
		DebitCost:  c.DebitCost,
		DebitOrg:   c.DebitOrg,
		CreditCost: c.CreditCost,
		CreditOrg:  c.CreditOrg,
	}
}

type EntryRequest struct {
	EntryTypeID   int64              `json:"entry_type_id" validate:"required"`
	Flow          string             `json:"flow" validate:"required"`
	Description   string             `json:"description" validate:"required"`
	AmountName    string             `json:"amount_name,omitempty"`
	Expression    string             `json:"expression,omitempty"`
	CashInBucket  string             `json:"cashin_bucket,omitempty"`
	CashOutBucket string             `json:"cashout_bucket,omitempty"`
	CostCenter    *CostCenterRequest `json:"cost_center,omitempty"`
	DebitAccount  *AccountRequest    `json:"debit_account,omitempty"`
	CreditAccount *AccountRequest    `json:"credit_account,omitempty"`
	Parameter     *ParameterRequest  `json:"parameter,omitempty"`
}

func (e EntryRequest) ToEntity() accounting.Entry {
	entry := accounting.Entry{
		EntryTypeID:   e.EntryTypeID,
		Flow:          e.Flow,
		Description:   e.Description,
		AmountName:    e.AmountName,
		Expression:    e.Expression,
		CashInBucket:  e.CashInBucket,
		CashOutBucket: e.CashOutBucket,
	}

	if e.CostCenter != nil {
		entry.CostCenter = e.CostCenter.ToEntity()
	}

	if e.DebitAccount != nil {
		entry.DebitAccount = e.DebitAccount.ToEntity()
	}

	if e.CreditAccount != nil {
		entry.CreditAccount = e.CreditAccount.ToEntity()
	}

	if e.Parameter != nil {
		entry.Parameter = e.Parameter.ToEntity()
	}

	return entry
}

func (e EntryRequest) Validate() error {
	if e.Flow == "" {
		return errors.New("entry.flow is required, choose an option: regular, migration")
	}

	if e.AmountName == "" && e.Expression == "" {
		return errors.New("entry.amount_name or expression is required, choose an option")
	}

	if e.AmountName != "" {
		if !slices.Contains(validAmounts, e.AmountName) {
			return fmt.Errorf("invalid entry.amount_name: %s", e.AmountName)
		}
	}

	if e.CreditAccount != nil {
		return e.CreditAccount.Validate()
	}

	if e.DebitAccount != nil {
		return e.DebitAccount.Validate()
	}

	if e.CostCenter != nil {
		return e.CostCenter.Validate()
	}

	if e.Parameter != nil {
		if err := e.Parameter.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type ParameterRequest struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func (p ParameterRequest) Validate() error {
	if p.Name == "" || p.Value == "" {
		return errors.New("entry.parameter.name and value are required")
	}
	return nil
}

func (p ParameterRequest) ToEntity() *accounting.Parameter {
	return &accounting.Parameter{
		Name:  p.Name,
		Value: p.Value,
	}
}
