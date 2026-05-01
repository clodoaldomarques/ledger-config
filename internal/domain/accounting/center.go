package accounting

type CostCenter struct {
	DebitCost  string
	DebitOrg   string
	CreditCost string
	CreditOrg  string
}

func (c CostCenter) Validate() error {
	return nil
}
