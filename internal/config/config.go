package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	DB       string `env:"DB_NAME" env-default:"diagram_editor"`
}

type Config struct {
	DBConfig
	Port string `env:"HTTP_PORT" env-default:"3000"`
}

func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return nil
	}
	return config
}

func Load() (*Config, error) {
	cfg := Config{}

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DB)
}
