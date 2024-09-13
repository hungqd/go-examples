package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbConnectionString string
}

func GetConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		DbConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}, nil
}

func MustGetConfig() *Config {
	config, err := GetConfig()
	if err != nil {
		log.Panicf("Get config error:  %v\n", err)
	}
	return config
}
