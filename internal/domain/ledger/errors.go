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
