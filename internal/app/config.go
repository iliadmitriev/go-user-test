package app

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen       string        `yaml:"listen" env:"LISTEN" env-defautl:":8000"`
	StoragePath  string        `yaml:"storage_path" env:"STORAGE_PATH" env-defautl:"main.db"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-defautl:"15s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-defautl:"15s"`
}

func NewConfig() (*Config, error) {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath != "" {
		cfgPath = "config.yaml"
	}

	if _, err := os.Stat(cfgPath); err != nil {
		return nil, err
	}

	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
