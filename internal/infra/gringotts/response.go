package gringotts

type GriphookResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}
