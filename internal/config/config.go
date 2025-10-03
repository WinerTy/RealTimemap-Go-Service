package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	user     string `yaml:"user" env-default:"postgres"`
	password string `yaml:"password" env-default:"postgres"`
	host     string `yaml:"host" env-default:"localhost"`
	port     int    `yaml:"port" env-default:"5432"`
	DbName   string `yaml:"db_name" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
type Config struct {
	Env        string `yaml:"env" env-default:"production"`
	Database   `yaml:"database"`
	HTTPServer `yaml:"http_server"`
}

func MustLoad() *Config {
	configPath := "./config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Cant read config file")
	}

	return &cfg
}
