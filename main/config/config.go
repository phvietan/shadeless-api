package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	bindAddress string
	frontendUrl string

	databaseUrl string

	// databaseHostname string
	// databaseUser     string
	// databasePassword string
	// databaseDbName   string
}

func (conf *config) init() *config {
	fmt.Println("Loading .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conf.bindAddress = os.Getenv("BIND_ADDRESS")
	conf.frontendUrl = os.Getenv("FRONTEND_URL")
	// conf.databaseUser = os.Getenv("DATABASE_USER")
	// conf.databaseDbName = os.Getenv("DATABASE_DBNAME")
	// conf.databaseHostname = os.Getenv("DATABASE_HOSTNAME")
	// conf.databasePassword = os.Getenv("DATABASE_PASSWORD")
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

// func (conf *config) GetDatabaseUser() string {
// 	return conf.databaseUser
// }
// func (conf *config) GetDatabaseDbName() string {
// 	return conf.databaseDbName
// }
// func (conf *config) GetDatabaseHostname() string {
// 	return conf.databaseHostname
// }
// func (conf *config) GetDatabasePassword() string {
// 	return conf.databasePassword
// }

var configInstance *config = new(config).init()

func GetInstance() *config { return configInstance }
