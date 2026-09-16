package config

import (
	"errors"
	"strings"

	"kialkuz/shop-with-loyalty/internal/config/db"
)

type Config struct {
	ServerHost           string
	ServerPort           string
	DB                   db.Config
	AccrualSystemAddress string
	SecretKey            string
}

func NewConfig() (*Config, error) {
	config, err := GetIncomingParams()
	if err != nil {
		return nil, err
	}

	if config.SecretKey == "" {
		return nil, errors.New("empty secret key")
	}

	addressParts := strings.Split(config.RunAddress, ":")

	return &Config{
		ServerHost:           addressParts[0],
		ServerPort:           addressParts[1],
		DB:                   *db.NewConfig(config.DatabaseURI),
		AccrualSystemAddress: config.AccrualSystemAddress,
		SecretKey:            config.SecretKey,
	}, nil
}
