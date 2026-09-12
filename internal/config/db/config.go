package db

type Config struct {
	DatabaseURI string
}

func NewConfig(databaseDSN string) *Config {
	return &Config{
		DatabaseURI: databaseDSN,
	}
}
