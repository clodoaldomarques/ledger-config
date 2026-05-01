package ledger

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
	CostCenter    *CostCenter
	DebitAccount  *Account
	CreditAccount *Account
}

func (e Script) Validate() error {
	return nil
}
