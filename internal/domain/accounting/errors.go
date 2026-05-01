package accounting

type ErrDuplicatedEntry struct {
	msg string
}

func (e ErrDuplicatedEntry) Error() string {
	if e.msg == "" {
		return "duplicated entry"
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
	return "accounting script not found"
}
