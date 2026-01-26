package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App App
	Db  Database
}

type App struct {
	Host    string `yaml:"host" env:"APP_HOST" env-default:"localhost"`
	Port    uint16 `yaml:"port" env:"APP_PORT" env-default:"8080"`
	Name    string `yaml:"name" env:"APP_NAME" env-default:"pizzeria"`
	Prefork bool   `yaml:"prefork" env:"APP_PREFORK" env-default:"false"`
}

type Database struct {
	Name string `yaml:"name" env:"DB_NAME" env-default:"pizzeria"`
	Host string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port uint16 `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Pass string `yaml:"pass" env:"DB_PASS" env-required:"true"`
}

var cfg Config

func New() (Config, error) {
	err := cleanenv.ReadEnv(&cfg)
	return cfg, err
}
