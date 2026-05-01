package dynamodb

type ErrScriptNotFound struct {
}

func (e ErrScriptNotFound) Error() string {
	return "script not found"
}
