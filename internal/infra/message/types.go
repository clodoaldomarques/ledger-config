package message

import (
	"time"

	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
)

type Config struct {
	ConfigID       string    `json:"config_id"`
	Level          string    `json:"event_id"`
	ProcessingCode string    `json:"processing_code"`
	OrgID          string    `json:"org_id"`
	ProgramID      int64     `json:"program_id"`
	Description    string    `json:"description"`
	Scripts        []Script  `json:"scripts"`
	Enable         bool      `json:"enable"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Version        int64     `json:"version"`
}

type Script struct {
	ScriptID      int64    `json:"script_id"`
	Flow          string   `json:"flow"`
	Description   string   `json:"description"`
	Expression    string   `json:"expression"`
	DebitAccount  *Account `json:"debit_account,omitempty"`
	CreditAccount *Account `json:"credit_account,omitempty"`
}

type Account struct {
	Number      string `json:"number"`
	Description string `json:"description"`
	Cosif       string `json:"cosif"`
}

func ToConfigMessage(c ledger.Config) Config {
	return Config{
		ConfigID:       c.ConfigID,
		Level:          string(c.Level),
		ProcessingCode: c.ProcessingCode,
		OrgID:          c.OrgID,
		ProgramID:      c.ProgramID,
		Description:    c.Description,
		Scripts:        ToScriptMessage(c.Scripts),
		Enable:         c.Enable,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
		Version:        c.Version,
	}
}

func ToScriptMessage(scs []ledger.Script) []Script {
	scm := make([]Script, len(scs))
	for _, sc := range scs {
		s := Script{
			ScriptID:      sc.ScriptID,
			Flow:          string(sc.Flow),
			Description:   sc.Description,
			Expression:    sc.Expression,
			DebitAccount:  ToAccountMessage(sc.DebitAccount),
			CreditAccount: ToAccountMessage(sc.CreditAccount),
		}
		scm = append(scm, s)
	}
	return scm
}

func ToAccountMessage(a *ledger.Account) *Account {
	if a != nil {
		return &Account{
			Number:      a.Number,
			Description: a.Description,
			Cosif:       a.Cosif,
		}
	}
	return nil
}
