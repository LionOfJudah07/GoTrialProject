package models

import "time"

type Transaction struct {
	ID        int64     `json:"id"`
	Amount    float64   `json:"amount"`
	Category  string    `json:"category"`
	Type      string    `json:"type"`
	Date      string    `json:"date"`
	Note      string    `json:"note"`
	AccountID int64     `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}