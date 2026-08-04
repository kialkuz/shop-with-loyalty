package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

const (
	defaultRunAddress           = "postgres"
	defaultDatabaseURI          = "localhost:8080"
	defaultAccrualSystemAddress = ""
)

type incomingParams struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func GetIncomingParams() (*incomingParams, error) {
	cfg := &incomingParams{
		RunAddress:           defaultRunAddress,
		DatabaseURI:          defaultDatabaseURI,
		AccrualSystemAddress: defaultAccrualSystemAddress,
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "Server address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Database address")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "Accrual address")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't parse config from os: %s ", err)
	}

	return cfg, nil
}
