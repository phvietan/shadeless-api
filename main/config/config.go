package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	environment string
	bindAddress string
	frontendUrl string
	databaseUrl string
}

func defaultValue(val string, defaultValue string) string {
	if val == "" {
		return defaultValue
	}
	return val
}

func (conf *config) init() *config {
	fmt.Println("Loading .env")
	_ = godotenv.Load()
	conf.environment = defaultValue(os.Getenv("ENVIRONMENT"), "test")
	conf.bindAddress = defaultValue(os.Getenv("BIND_ADDRESS"), "0.0.0.0:3000")
	conf.frontendUrl = defaultValue(os.Getenv("FRONTEND_URL"), "")
	conf.databaseUrl = defaultValue(os.Getenv("DATABASE_URL"), "")
	return conf
}

func (conf *config) GetEnvironment() string {
	return conf.environment
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

func (conf *config) SetEnvironment(val string) {
	conf.environment = val
}
func (conf *config) SetBindAddress(val string) {
	conf.bindAddress = val
}
func (conf *config) SetFrontendUrl(val string) {
	conf.frontendUrl = val
}
func (conf *config) SetDatabaseUrl(val string) {
	conf.databaseUrl = val
}

var configInstance *config = new(config).init()

func GetInstance() *config { return configInstance }
