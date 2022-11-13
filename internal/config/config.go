package config

import "os"

type (
	HTTPConfig struct {
	}

	PostgresConfig struct {
	}

	Config struct {
		HTTP     HTTPConfig
		Postgres PostgresConfig
	}
)

func Init(configDir) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config

	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func populateDefaults() {

}

func parseConfigFile(dir string, env string) error {
	return nil
}

func unmarshal(c *Config) error {
	return nil
}

func setFromEnv(c *Config) {

}
