package models

type Registration struct {
	ID  int64  `json:"id,omitempty"`
	CIF string `json:"cif,omitempty"`
}
