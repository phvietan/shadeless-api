package config

import (
	"shadeless-api/main/libs"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNonEnvFileInit(t *testing.T) {
	var conf = new(config).init()
	assert.Equal(t, conf.bindAddress, "0.0.0.0:3000")
	assert.Equal(t, conf.frontendUrl, "")
}

func TestInit(t *testing.T) {
	var conf = new(config).init()
	assert.Equal(t, conf.bindAddress, "0.0.0.0:3000")
	assert.Equal(t, conf.environment, "test")
	assert.Equal(t, conf.frontendUrl, "")
}

func TestGetInstance(t *testing.T) {
	c := GetInstance()
	assert.Equal(t, c.bindAddress, "0.0.0.0:3000")
	assert.Equal(t, c.environment, "test")
	assert.Equal(t, c.frontendUrl, "")
}

func TestSetInstance(t *testing.T) {
	c := GetInstance()
	environment := libs.RandomString(32)
	bindAddress := libs.RandomString(32)
	frontendUrl := libs.RandomString(32)
	databaseUrl := libs.RandomString(32)
	c.SetEnvironment(environment)
	c.SetBindAddress(bindAddress)
	c.SetFrontendUrl(frontendUrl)
	c.SetDatabaseUrl(databaseUrl)

	assert.Equal(t, c.GetEnvironment(), environment)
	assert.Equal(t, c.GetBindAddress(), bindAddress)
	assert.Equal(t, c.GetFrontendUrl(), frontendUrl)
	assert.Equal(t, c.GetDatabaseUrl(), databaseUrl)
}
