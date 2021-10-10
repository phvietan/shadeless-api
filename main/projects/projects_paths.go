package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
)

func PathsRoutes(route *gin.Engine) {
	users := route.Group("/projects/:projectName")
	{
		users.GET("/paths", getPathsByOrigin)
	}
}

func getPathsByOrigin(c *gin.Context) {
	projectName := c.Param("projectName")
	origin := c.Query("origin")

	var pathDb database.IParsedPathDatabase = new(database.ParsedPathDatabase).Init()
	paths := pathDb.GetPathsByProjectAndOrigin(projectName, origin)
	responser.ResponseOk(c, paths)
}
