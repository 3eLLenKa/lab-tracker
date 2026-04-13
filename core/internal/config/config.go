package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env    string `yaml:"env" env:"ENV" env-default:"dev"`
	Server Server `yaml:"server"`
	DB     DB     `yaml:"db"`
}

type Server struct {
	Host        string        `yaml:"host" env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port        string        `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

type DB struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"POSTGRES_DB" env-default:"postgres"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"postgres"`
	SSLMode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" env-default:"disable"`
}

func MustLoad() *Config {
	for _, path := range []string{".env", "../.env"} {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Overload(path); err != nil {
				log.Printf("failed to load %s: %v", path, err)
			}
		}
	}

	cfg := &Config{}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Printf("failed to read env: %v", err)
	}

	return cfg
}
