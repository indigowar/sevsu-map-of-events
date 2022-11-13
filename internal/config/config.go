package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegaBytes = 1

	envLocal = "local"
)

type (
	HTTPConfig struct {
		Port                string        `mapstructure:"port"`
		ReadTimeout         time.Duration `mapstructure:"readTimeout"`
		WriteTimeout        time.Duration `mapstructure:"writeTimeout"`
		MaxHeadersMegabytes int           `mapstructure:"maxHeaderMegaBytes"`
	}

	PostgresConfig struct {
	}

	Config struct {
		HTTP        HTTPConfig
		Postgres    PostgresConfig
		Environment string
	}
)

func Init(configDir string) (*Config, error) {
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
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegaBytes)
}

func parseConfigFile(dir string, env string) error {
	viper.AddConfigPath(dir)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == envLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func unmarshal(c *Config) error {
	if err := viper.UnmarshalKey("http", &c.HTTP); err != nil {
		return err
	}

	return nil
}

func setFromEnv(c *Config) {
	c.Environment = os.Getenv("APP_ENV")
}
