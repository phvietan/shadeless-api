package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type config struct {
	bindAddress string
	databaseUrl string
	frontendUrl string
}

func (conf *config) init() *config {
	fmt.Println("Loading .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conf.databaseUrl = os.Getenv("DATABASE_URL")
	conf.bindAddress = os.Getenv("BIND_ADDRESS")
	conf.frontendUrl = os.Getenv("FRONTEND_URL")
	return conf
}

func (conf *config) GetDatabaseURl() string {
	return conf.databaseUrl
}
func (conf *config) GetBindAddress() string {
	return conf.bindAddress
}
func (conf *config) GetFrontendUrl() string {
	return conf.frontendUrl
}

var lock = &sync.Mutex{}
var configInstance *config

func GetInstance() *config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if configInstance == nil {
			configInstance = new(config).init()
		} else {
			fmt.Println("Single Instance already created-1")
		}
	} else {
		fmt.Println("Single Instance already created-2")
	}
	return configInstance
}
