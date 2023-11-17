package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	HttpServer
	Database
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"Abdrahman"`
	DBname   string `yaml:"dbname" env-default:"gowebsocket"`
	Hostname string `yaml:"hostname" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
}

func MustLoad() Config {
	os.Setenv("CONFIG_PATH", "../../config/local.yaml")

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config from file: %s", err)
	}

	return cfg
}
