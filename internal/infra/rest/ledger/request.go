package ledger

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
)

type PostConfigRequest struct {
	Level       string           `json:"level" validate:"required"`
	ProcessCode string           `json:"process_code" validate:"required"`
	ProgramID   *int64           `json:"program_id,omitempty"`
	Description string           `json:"description" validate:"required"`
	Accounts    []AccountRequest `json:"accounts,omitempty"`
	Scripts     []ScriptRequest  `json:"scripts" validate:"required"`
	Enable      *bool            `json:"enable,omitempty"`
}

func (p PostConfigRequest) Validate() error {
	if p.Level == "" {
		return fmt.Errorf("level is required")
	}

	if p.Description == "" {
		return fmt.Errorf("description is required")
	}

	if len(p.Scripts) == 0 {
		return fmt.Errorf("scripts is required")
	}

	for _, e := range p.Scripts {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (p PostConfigRequest) PostToEntity(orgID string) ledger.Config {
	scr := ledger.Config{
		Level:       ledger.Level(p.Level),
		ProcessCode: strings.ToUpper(p.ProcessCode),
		OrgID:       orgID,
		Description: p.Description,
		Scripts:     make([]ledger.Script, 0, len(p.Scripts)),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Enable:      true,
		Version:     1,
	}

	if p.ProgramID != nil {
		scr.ProgramID = *p.ProgramID
	}

	for _, e := range p.Scripts {
		scr.Scripts = append(scr.Scripts, e.ToEntity())
	}

	return scr
}

type PathScriptRequest struct {
	Description string          `json:"description" validate:"required"`
	Scripts     []ScriptRequest `json:"scripts" validate:"required"`
	Enable      *bool           `json:"enable,omitempty"`
}

func (p PathScriptRequest) Validate() error {

	if p.Description == "" {
		return fmt.Errorf("description is required")
	}

	if len(p.Scripts) == 0 {
		return fmt.Errorf("scripts is required")
	}

	for _, s := range p.Scripts {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (p PathScriptRequest) PatchToEntity(orgID string) ledger.Config {
	scr := ledger.Config{
		OrgID:       orgID,
		Description: p.Description,
		Scripts:     make([]ledger.Script, 0, len(p.Scripts)),
	}

	for _, s := range p.Scripts {
		scr.Scripts = append(scr.Scripts, s.ToEntity())
	}

	if p.Enable != nil {
		scr.Enable = *p.Enable
	}

	return scr
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

func (a AccountRequest) ToEntity() *ledger.Account {
	return &ledger.Account{
		Number:      a.Number,
		Description: a.Description,
		Cosif:       a.Cosif,
	}
}

type ScriptRequest struct {
	ScriptID      int64           `json:"script_id" validate:"required"`
	Flow          string          `json:"flow" validate:"required"`
	Description   string          `json:"description" validate:"required"`
	AmountName    string          `json:"amount_name,omitempty"`
	Expression    string          `json:"expression,omitempty"`
	DebitAccount  *AccountRequest `json:"debit_account,omitempty"`
	CreditAccount *AccountRequest `json:"credit_account,omitempty"`
}

func (e ScriptRequest) ToEntity() ledger.Script {
	entry := ledger.Script{
		ScriptID:    e.ScriptID,
		Flow:        e.Flow,
		Description: e.Description,
		Expression:  e.Expression,
	}

	if e.DebitAccount != nil {
		entry.DebitAccount = e.DebitAccount.ToEntity()
	}

	if e.CreditAccount != nil {
		entry.CreditAccount = e.CreditAccount.ToEntity()
	}

	return entry
}

func (e ScriptRequest) Validate() error {
	if e.Flow == "" {
		return errors.New("script.flow is required, choose an option: regular, migration")
	}

	if e.Expression == "" {
		return errors.New("scripts.expression is required")
	}

	if e.CreditAccount != nil {
		return e.CreditAccount.Validate()
	}

	if e.DebitAccount != nil {
		return e.DebitAccount.Validate()
	}

	return nil
}
