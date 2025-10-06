package config

import (
	"fmt"
	"log"
	"os"
	"realtimemap-service/internal/pkg/logger/sl"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"env" env-default:"production"`
	CacheStrategy string `yaml:"cache_strategy" env-default:"noop"`
	Database      `yaml:"database"`
	HTTPServer    `yaml:"http_server"`
	Redis         `yaml:"redis"`
}

type Database struct {
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	DbName   string `yaml:"db_name" env-required:"true"`
}

type Redis struct {
	Url      string `yaml:"url" env-default:"localhost:6379"`
	Password string `yaml:"password" env-default:""`
	DB       int    `yaml:"db" env-default:"0"`
}

func (d *Database) BuildURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.DbName)
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := "./config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config file does not exist", sl.Err(err))
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Cant read config file", sl.Err(err))
	}

	return &cfg
}
