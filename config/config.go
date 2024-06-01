package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Conf struct {
	General  *ConfGeneral
	API      *ConfAPI
	Database *ConfDatabase
}

type ConfGeneral struct {
	Tag string
}

type ConfAPI struct {
	Port string
}

type ConfDatabase struct {
	Driver       string
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}

// TODO: Use struct tags to minimize manual effort (or just use an env management tool)
func New() *Conf {
	godotenv.Load()

	var config Conf
	config.General = &ConfGeneral{
		Tag: lookup("TAG"),
	}
	config.API = &ConfAPI{
		Port: lookup("API_PORT"),
	}
	config.Database = &ConfDatabase{
		Driver:       lookup("DATABASE_DRIVER"),
		Host:         lookup("DATABASE_HOST"),
		Port:         lookup("DATABASE_PORT"),
		Username:     lookup("DATABASE_USER"),
		Password:     lookup("DATABASE_PASS"),
		DatabaseName: lookup("DATABASE_NAME"),
		SSLMode:      lookup("DATABASE_SSL"),
	}
	return &config
}

func lookup(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Environment variable `%s` is missing\n", key)
	}
	return val
}
