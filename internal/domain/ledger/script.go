package ledger

import "fmt"

type Flow = string

const (
	Regular   Flow = "regular"
	Migration Flow = "migration"
)

type Script struct {
	ScriptID      int64
	Flow          Flow
	Description   string
	Expression    string
	DebitAccount  *Account
	CreditAccount *Account
}

func (s Script) ScriptKey() string {
	return fmt.Sprintf("%1s#%2d", s.Flow, s.ScriptID)
}

func (s Script) Validate() error {
	return nil
}
