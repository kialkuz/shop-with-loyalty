package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

const (
	defaultRunAddress           = "localhost:8081"
	defaultDatabaseURI          = ""
	defaultAccrualSystemAddress = ""
	defaultSecretKey            = ""
)

type incomingParams struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey            string `env:"SECRET_KEY"`
}

func GetIncomingParams() (*incomingParams, error) {
	cfg := &incomingParams{
		RunAddress:           defaultRunAddress,
		DatabaseURI:          defaultDatabaseURI,
		AccrualSystemAddress: defaultAccrualSystemAddress,
		SecretKey:            defaultSecretKey,
	}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't parse config from os: %s ", err)
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "Server address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Database address")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "Accrual address")
	flag.StringVar(&cfg.SecretKey, "k", cfg.SecretKey, "Secret key")
	flag.Parse()

	return cfg, nil
}
