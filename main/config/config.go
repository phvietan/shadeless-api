package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	bindAddress string
	frontendUrl string
	databaseUrl string
}

func (conf *config) init() *config {
	fmt.Println("Loading .env")
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return nil
	}

	conf.bindAddress = os.Getenv("BIND_ADDRESS")
	conf.frontendUrl = os.Getenv("FRONTEND_URL")
	conf.databaseUrl = os.Getenv("DATABASE_URL")
	return conf
}

func (conf *config) GetBindAddress() string {
	return conf.bindAddress
}
func (conf *config) GetFrontendUrl() string {
	return conf.frontendUrl
}
func (conf *config) GetDatabaseUrl() string {
	return conf.databaseUrl
}

var configInstance *config = new(config).init()

func GetInstance() *config { return configInstance }
