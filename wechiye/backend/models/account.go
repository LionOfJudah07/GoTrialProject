package models

import "time"

type Account struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	InitialBalance float64   `json:"initial_balance"`
	CurrentBalance float64   `json:"current_balance"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}