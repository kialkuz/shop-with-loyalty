package config

import (
	"strings"

	"kialkuz/shop-with-loyalty/internal/config/db"
)

type Config struct {
	ServerHost        string
	ServerPort        string
	DB                db.Config
	AccrualServerHost string
	AccrualServerPort string
	SecretKey         string
}

func NewConfig() (*Config, error) {
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
		SecretKey:         config.SecretKey,
	}, nil
}
