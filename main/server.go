package main

import (
	"fmt"
	"shadeless-api/main/burp"
	"shadeless-api/main/config"
	"shadeless-api/main/projects"
	"strings"

	"github.com/gin-gonic/gin"
)

func healthCheckHandler(c *gin.Context) {
	c.String(200, "Health check ok")
}

func spawnApp() *gin.Engine {
	router := gin.Default()

	router.Use(setHeaderOctetStream())
	router.Static("/files", "./files")

	router.Use(setHeaderForApi())
	router.Use(handleOptionsMethod())
	router.GET("/healthcheck", healthCheckHandler)

	burp.Routes(router)
	projects.Routes(router)
	return router
}

func main() {
	s := "/"
	a := strings.Split(s, "/")
	fmt.Println(a, len(a))

	router := spawnApp()
	router.Run(config.GetInstance().GetBindAddress()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
