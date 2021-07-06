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

func main() {
	router := gin.Default()
	router.GET("/healthcheck", healthCheckHandler)

	router.Use(setHeaderOctetStream())
	router.Static("/files", "./files")

	router.Use(setHeaderForApi())
	router.Use(handleOptionsMethod())

	burp.Routes(router)
	projects.Routes(router)
	router.Run(config.GetInstance().GetBindAddress()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
