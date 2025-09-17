package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`
	Pg  `yaml:"postgres" env-required:"true"`
}

type Pg struct {
	Url string `yaml:"url" env-required:"true"`
}

func MustLoad() *Config {
	configPathFromEnv := os.Getenv("CONFIG_PATH")
	if configPathFromEnv == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	return MustLoadByPath(configPathFromEnv)
}


// Loads the config structure from the file path
func MustLoadByPath(pathToConfigFile string) *Config {
	if _, err := os.Stat(pathToConfigFile); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", pathToConfigFile)
	}

	var config Config

	if err := cleanenv.ReadConfig(pathToConfigFile, &config); err != nil {
		log.Fatal(err)
	}

	return &config
}
