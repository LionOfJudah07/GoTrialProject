package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/curve25519"
)

func GenerateX25519KeyPair() (pubBase64, privBase64 string, err error) {
	var priv, pub [32]byte
	if _, err := rand.Read(priv[:]); err != nil {
		return "", "", err
	}
	priv[0] &= 248
	priv[31] &= 127
	priv[31] |= 64

	curve25519.ScalarBaseMult(&pub, &priv)
	return base64.StdEncoding.EncodeToString(pub[:]), base64.StdEncoding.EncodeToString(priv[:]), nil
}

func DeriveSharedSecret(privBase64, peerPubBase64 string) (string, error) {
	privBytes, err := base64.StdEncoding.DecodeString(privBase64)
	if err != nil || len(privBytes) != 32 {
		return "", errors.New("invalid private key")
	}
	pubBytes, err := base64.StdEncoding.DecodeString(peerPubBase64)
	if err != nil || len(pubBytes) != 32 {
		return "", errors.New("invalid public key")
	}
	var secret [32]byte
	curve25519.ScalarMult(&secret, (*[32]byte)(privBytes), (*[32]byte)(pubBytes))
	return base64.StdEncoding.EncodeToString(secret[:]), nil
}