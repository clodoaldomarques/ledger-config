package ledger

type Flow = string

const (
	Regular   Flow = "regular"
	Migration Flow = "migration"
)

type Entry struct {
	EntryTypeID   int64
	Flow          Flow
	Description   string
	AmountName    string
	Expression    string
	CashInBucket  string
	CashOutBucket string
	CostCenter    *CostCenter
	DebitAccount  *Account
	CreditAccount *Account
	Parameter     *Parameter
}

func (e Entry) Validate() error {
	return nil
}
