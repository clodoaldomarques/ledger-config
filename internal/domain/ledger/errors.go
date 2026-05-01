package ledger

type ErrDuplicatedScript struct {
	msg string
}

func (e ErrDuplicatedScript) Error() string {
	if e.msg == "" {
		return "duplicated script"
	}
	return e.msg
}

type ErrDuplicatedAccount struct {
	msg string
}

func (e ErrDuplicatedAccount) Error() string {
	if e.msg == "" {
		return "duplicated account type"
	}
	return e.msg
}

type ErrScriptNotFound struct {
}

func (e ErrScriptNotFound) Error() string {
	return "ledger config not found"
}
