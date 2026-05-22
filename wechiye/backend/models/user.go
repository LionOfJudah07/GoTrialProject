package models

import "time"

type User struct {
	ID                    int64     `json:"id"`
	FullName              string    `json:"full_name"`
	Email                 string    `json:"email"`
	Username              string    `json:"username"`
	Avatar                string    `json:"avatar"`
	Gender                string    `json:"gender"`
	EducationLevel        string    `json:"education_level"`
	Occupation            string    `json:"occupation"`
	HasKids               bool      `json:"has_kids"`
	KidsAllowanceAmount   float64   `json:"kids_allowance_amount"`
	KidsAllowanceInterval string    `json:"kids_allowance_interval"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}