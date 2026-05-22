package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"wechiye/backend/crypto"
	"wechiye/backend/database"
	"wechiye/backend/models"
)

type CoupleService struct {
	db *database.DB
}

func NewCoupleService(db *database.DB) *CoupleService {
	return &CoupleService{db: db}
}

func (s *CoupleService) GetLocalPublicKey() (string, error) {
	data, err := s.getOrCreateCoupleData()
	if err != nil {
		return "", err
	}
	return data.PublicKey, nil
}

func (s *CoupleService) getOrCreateCoupleData() (*models.CoupleData, error) {
	row := s.db.QueryRow(`SELECT id, device_id, public_key, private_key, partner_public_key, shared_secret, linked_at, created_at, updated_at 
		FROM couple_data LIMIT 1`)
	var data models.CoupleData
	var linkedAt *time.Time
	err := row.Scan(&data.ID, &data.DeviceID, &data.PublicKey, &data.PrivateKey,
		&data.PartnerPublicKey, &data.SharedSecret, &linkedAt, &data.CreatedAt, &data.UpdatedAt)
	if err == nil {
		data.LinkedAt = linkedAt
		return &data, nil
	}
	deviceID := generateDeviceID()
	pub, priv, err := crypto.GenerateX25519KeyPair()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_, err = s.db.Exec(`INSERT INTO couple_data (device_id, public_key, private_key, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`, deviceID, pub, priv, now, now)
	if err != nil {
		return nil, err
	}
	return &models.CoupleData{
		DeviceID:  deviceID,
		PublicKey: pub,
		PrivateKey: priv,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func generateDeviceID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (s *CoupleService) GenerateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (s *CoupleService) ProcessPairingQR(qrData string) error {
	if !strings.HasPrefix(qrData, "wechiye:link:") {
		return errors.New("invalid QR code")
	}
	parts := strings.Split(qrData, ":")
	if len(parts) != 5 {
		return errors.New("invalid QR data")
	}
	partnerPub := parts[2]
	version := parts[4]
	if version != "1.0" {
		return fmt.Errorf("unsupported version %s", version)
	}
	data, err := s.getOrCreateCoupleData()
	if err != nil {
		return err
	}
	sharedSecret, err := crypto.DeriveSharedSecret(data.PrivateKey, partnerPub)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = s.db.Exec(`UPDATE couple_data SET partner_public_key=?, shared_secret=?, linked_at=?, updated_at=? WHERE id=?`,
		partnerPub, sharedSecret, now, now, data.ID)
	return err
}

func (s *CoupleService) GetStatus() (map[string]interface{}, error) {
	data, err := s.getOrCreateCoupleData()
	if err != nil {
		return nil, err
	}
	linked := data.PartnerPublicKey != ""
	status := map[string]interface{}{
		"linked":    linked,
		"device_id": data.DeviceID,
	}
	if linked && data.LinkedAt != nil {
		status["linked_at"] = data.LinkedAt
	}
	return status, nil
}

func (s *CoupleService) GetSharedSecret() ([]byte, error) {
	data, err := s.getOrCreateCoupleData()
	if err != nil {
		return nil, err
	}
	if data.SharedSecret == "" {
		return nil, errors.New("not linked")
	}
	return base64.StdEncoding.DecodeString(data.SharedSecret)
}