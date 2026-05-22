package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "wechiye"
	keyName     = "master-key"
)

type Config struct {
	DataDir string
}

func NewConfig(dataDir string) *Config {
	return &Config{DataDir: dataDir}
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	var dir string
	if os.Getenv("APPDATA") != "" {
		dir = filepath.Join(os.Getenv("APPDATA"), "wechiye")
	} else if _, err := os.Stat("/Library"); err == nil {
		dir = filepath.Join(home, "Library", "Application Support", "wechiye")
	} else {
		dir = filepath.Join(home, ".wechiye")
	}
	return dir, nil
}

func (c *Config) SaveMasterKey(key []byte) error {
	hexKey := fmt.Sprintf("%x", key)
	return keyring.Set(serviceName, keyName, hexKey)
}

func (c *Config) LoadMasterKey() ([]byte, error) {
	hexKey, err := keyring.Get(serviceName, keyName)
	if err != nil {
		return nil, err
	}
	var key []byte
	_, err = fmt.Sscanf(hexKey, "%x", &key)
	return key, err
}

func (c *Config) DeleteMasterKey() error {
	return keyring.Delete(serviceName, keyName)
}