package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Addr      string `yaml:"addr"`
	DBurl     string `yaml:"dburl"`
	JWTSecret string `yaml:"jwt_secret"`
}

func Load(path string) *Config {
	var conf Config
	if err := cleanenv.ReadConfig(path, &conf); err != nil {
		log.Fatal("couldn't read config")
	}
	return &conf
}
