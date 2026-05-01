package ledger

type Company struct {
	Code string
	Type string
}

func (c Company) Validate() error {
	return nil
}
