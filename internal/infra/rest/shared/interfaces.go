package shared

type EntityRequest interface {
	Validate() error
}
