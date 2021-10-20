package main

import (
	"net/http"
	"shadeless-api/main/config"

	"github.com/gin-gonic/gin"
)

func setHeaderOctetStream() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Access-Control-Allow-Origin", config.GetInstance().GetFrontendUrl())
		c.Next()
	}
}

func setHeaderForApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")
		c.Next()
	}
}

func handleOptionsMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
