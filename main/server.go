package main

import (
	"shadeless-api/main/burp"
	"shadeless-api/main/config"
	"shadeless-api/main/projects"

	"github.com/gin-gonic/gin"
)

func healthCheckHandler(c *gin.Context) {
	c.String(200, "Health check ok")
}

func spawnApp() *gin.Engine {
	router := gin.Default()

	router.Use(setHeaderOctetStream())
	router.Static("/files", "../files")
	router.Use(setHeaderForApi())
	router.Use(handleOptionsMethod())

	router.GET("/healthcheck", healthCheckHandler)

	burp.Routes(router)
	projects.Routes(router)
	return router
}

func main() {
	router := spawnApp()
	router.Run(config.GetInstance().GetBindAddress())
}
