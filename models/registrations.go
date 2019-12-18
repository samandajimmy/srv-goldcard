package models

// Registrations is a struct to store registration data
type Registrations struct {
	ID  int64  `json:"id,omitempty"`
	CIF string `json:"cif,omitempty"`
}
