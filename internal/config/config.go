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
	Port int    `koanf:"port" validate:"required"`
	Env  string `koanf:"env" validate:"required"`
}

type DBConfig struct {
	URL             string `koanf:"url" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns"`
	MaxIdleConns    int    `koanf:"max_idle_conns"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime"`
}

type AuthConfig struct {
	JWTSecret     string `koanf:"jwt_secret" validate:"required"`
	TokenDuration int    `koanf:"token_duration"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	k := koanf.New(".")

	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}), nil); err != nil {
		return nil, err
	}

	setDefaults(k)

	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults(k *koanf.Koanf) {
	k.Set("app.port", 8080)
	k.Set("app.env", "dev")

	k.Set("db.max_open_conns", 25)
	k.Set("db.max_idle_conns", 10)
	k.Set("db.conn_max_lifetime", 300)

	k.Set("auth.token_duration", 60)
}
