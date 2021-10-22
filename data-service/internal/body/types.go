package body

type body struct {
	IsOk    bool        `json:"is_ok"`
	Payload interface{} `json:"payload,omitempty"`
}
