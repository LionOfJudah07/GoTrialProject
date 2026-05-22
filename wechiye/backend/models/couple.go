package models

import "time"

type CoupleData struct {
	ID               int64      `json:"id"`
	DeviceID         string     `json:"device_id"`
	PublicKey        string     `json:"public_key"`
	PrivateKey       string     `json:"private_key"`
	PartnerPublicKey string     `json:"partner_public_key"`
	SharedSecret     string     `json:"shared_secret"`
	LinkedAt         *time.Time `json:"linked_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}