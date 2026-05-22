package database

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"

	"golang.org/x/crypto/scrypt"
)

const (
	scryptN      = 32768
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 32
)

func DeriveKey(password string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	key, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptKeyLen)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func ZeroBytes(b []byte) {
	subtle.ConstantTimeCopy(1, b, make([]byte, len(b)))
}

func BytesToHexString(b []byte) string {
	return hex.EncodeToString(b)
}