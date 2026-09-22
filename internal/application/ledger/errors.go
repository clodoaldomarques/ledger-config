package ledger

import "fmt"

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

type ErrConfigNotFound struct {
}

func (e ErrConfigNotFound) Error() string {
	return "ledger config not found"
}

type ErrOrgActivated struct {
	OrgID string
}

func (e ErrOrgActivated) Error() string {
	return fmt.Sprintf("tenant %s was activated", e.OrgID)
}
