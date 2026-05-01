package shared

type ErrResponse struct {
	Message string `json:"message"`
}

func (e ErrResponse) Error() string {
	return e.Message
}
