package config

import (
	"strings"

	"kialkuz/gophermart/internal/config/db"
	"kialkuz/gophermart/internal/infrastructure/env"
)

type Config struct {
	ServerHost        string
	ServerPort        string
	DB                db.Config
	AccrualServerHost string
	AccrualServerPort string
}

func NewConfig() (*Config, error) {
	err := env.Load()
	if err != nil {
		return nil, err
	}

	config, err := GetIncomingParams()
	if err != nil {
		return nil, err
	}

	addressParts := strings.Split(config.RunAddress, ":")
	accrualAddressParts := strings.Split(config.AccrualSystemAddress, ":")

	return &Config{
		ServerHost:        addressParts[0],
		ServerPort:        addressParts[1],
		DB:                *db.NewConfig(config.DatabaseURI),
		AccrualServerHost: accrualAddressParts[0],
		AccrualServerPort: accrualAddressParts[1],
	}, nil
}
