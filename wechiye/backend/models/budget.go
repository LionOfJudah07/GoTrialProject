package models

import "time"

type Budget struct {
	ID          int64     `json:"id"`
	Category    string    `json:"category"`
	MonthYear   string    `json:"month_year"`
	AmountLimit float64   `json:"amount_limit"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}