package config

import (
	"io/ioutil"
	"log"
	"os"
	"shadeless-api/main/libs"
	"testing"

	"github.com/go-playground/assert/v2"
)

func generateRandomEnvFile() (string, string, string) {
	databaseUrl := libs.RandomString(32)
	bindAddress := libs.RandomString(32)
	frontendUrl := libs.RandomString(32)
	content := "DATABASE_URL=" + databaseUrl + "\n" + "BIND_ADDRESS=" + bindAddress + "\n" + "FRONTEND_URL=" + frontendUrl + "\n"
	if err := ioutil.WriteFile(".env", []byte(content), 0755); err != nil {
		log.Fatal("Unable to write file")
	}
	return databaseUrl, bindAddress, frontendUrl
}

func removeEnvFile() {
	os.Remove(".env")
}

func TestNonEnvFileInit(t *testing.T) {
	var conf = new(config).init()
	assert.Equal(t, nil, conf)
	removeEnvFile()
}

func TestInit(t *testing.T) {
	databaseUrl, bindAddress, frontendUrl := generateRandomEnvFile()
	var conf = new(config).init()
	assert.Equal(t, conf.bindAddress, bindAddress)
	assert.Equal(t, conf.databaseUrl, databaseUrl)
	assert.Equal(t, conf.frontendUrl, frontendUrl)
	assert.Equal(t, conf.GetBindAddress(), bindAddress)
	assert.Equal(t, conf.GetDatabaseUrl(), databaseUrl)
	assert.Equal(t, conf.GetFrontendUrl(), frontendUrl)
	removeEnvFile()
}
