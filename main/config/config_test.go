package config

import (
	"io/ioutil"
	"log"
	"os"
	"shadeless-api/main/libs"
	"testing"

	"github.com/go-playground/assert/v2"
)

func generateRandomEnvFile() (string, string, string, string) {
	environment := libs.RandomString(32)
	databaseUrl := libs.RandomString(32)
	bindAddress := libs.RandomString(32)
	frontendUrl := libs.RandomString(32)
	content := "ENVIRONMENT=" + environment + "\nDATABASE_URL=" + databaseUrl + "\nBIND_ADDRESS=" + bindAddress + "\nFRONTEND_URL=" + frontendUrl
	if err := ioutil.WriteFile(".env", []byte(content), 0755); err != nil {
		log.Fatal("Unable to write file")
	}
	return environment, databaseUrl, bindAddress, frontendUrl
}

func removeEnvFile() {
	os.Remove(".env")
}

func TestNonEnvFileInit(t *testing.T) {
	var conf = new(config).init()
	assert.Equal(t, conf.bindAddress, "0.0.0.0:3000")
	assert.Equal(t, conf.frontendUrl, "")
}

func TestInit(t *testing.T) {
	removeEnvFile()
	environment, databaseUrl, bindAddress, frontendUrl := generateRandomEnvFile()
	var conf = new(config).init()
	assert.Equal(t, conf.environment, environment)
	assert.Equal(t, conf.bindAddress, bindAddress)
	assert.Equal(t, conf.databaseUrl, databaseUrl)
	assert.Equal(t, conf.frontendUrl, frontendUrl)
	assert.Equal(t, conf.GetEnvironment(), environment)
	assert.Equal(t, conf.GetBindAddress(), bindAddress)
	assert.Equal(t, conf.GetDatabaseUrl(), databaseUrl)
	assert.Equal(t, conf.GetFrontendUrl(), frontendUrl)
	removeEnvFile()
}

func TestGetInstance(t *testing.T) {
	removeEnvFile()
	c := GetInstance()
	assert.Equal(t, c.bindAddress, "0.0.0.0:3000")
	assert.Equal(t, c.environment, "test")
	assert.Equal(t, c.frontendUrl, "")
}
