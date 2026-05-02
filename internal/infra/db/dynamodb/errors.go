package dynamodb

type ErrConfigNotFound struct {
}

func (e ErrConfigNotFound) Error() string {
	return "config not found"
}
