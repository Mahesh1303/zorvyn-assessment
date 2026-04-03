package config

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	App  AppConfig  `koanf:"app"`
	DB   DBConfig   `koanf:"db"`
	Auth AuthConfig `koanf:"auth"`
}

type AppConfig struct {
	Port int `koanf:"port"`
}

type DBConfig struct {
	URL string `koanf:"url"`
}

type AuthConfig struct {
	JWTSecret string `koanf:"jwt_secret"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// creating Koanff instanceee so that we can load the env
	k := koanf.New(".")

	// Loading the env file
	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", 1)
	}), nil); err != nil {
		return nil, err
	}

	// unmarshalling the envs into the config Struct
	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	// validating the env types to its struct types
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
