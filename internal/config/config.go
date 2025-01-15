package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

type Config struct {
	HTTP     HttpConfig
	PG       PostgresConfig
	LogLevel string `env:"LOG_LEVEL" env-default:"debug"`
}

type HttpConfig struct {
	Host string `env:"HTTP_HOST" env-default:"localhost"`
	Port string `env:"HTTP_PORT" env-default:"8080"`
}

type PostgresConfig struct {
	DSN string `env:"POSTGRES_DSN"`
}

func MustNew() Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
