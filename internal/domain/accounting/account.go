package accounting

type Account struct {
	Number      string
	Description string
	Cosif       string
}

func (a Account) Validate() error {
	return nil
}
