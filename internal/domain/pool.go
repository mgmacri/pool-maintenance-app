package domain

import "time"

// PoolType represents the sanitization system.
type PoolType string

const (
	PoolTypeChlorine PoolType = "chlorine"
	PoolTypeSalt     PoolType = "salt"
)

// Pool represents a swimming pool to be maintained.
type Pool struct {
	ID        string    `json:"id" example:"pool-123"`
	Name      string    `json:"name" example:"Backyard Oasis"`
	Gallons   int       `json:"gallons" example:"15000"`
	Type      PoolType  `json:"type" example:"salt"`
	CreatedAt time.Time `json:"created_at"`
}

// ChemicalReading represents a snapshot of water chemistry.
type ChemicalReading struct {
	PoolID       string  `json:"pool_id" example:"pool-123"`
	PH           float64 `json:"ph" example:"7.2"`
	FreeChlorine float64 `json:"free_chlorine" example:"3.0"`
	Alkalinity   int     `json:"alkalinity" example:"80"`
	CyanuricAcid int     `json:"cyanuric_acid" example:"30"`
}

// JobReport represents the output of a maintenance visit.
type JobReport struct {
	PoolID         string   `json:"pool_id"`
	IsBalanced     bool     `json:"is_balanced"`
	Dosage         []string `json:"dosage" example:"Add 1lb Soda Ash"`
	TechnicianNote string   `json:"technician_note,omitempty"`
}
