package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
)

func PathsRoutes(route *gin.Engine) {
	parsedPathRoute := route.Group("/projects/:projectName")
	{
		parsedPathRoute.GET("/paths", getPathsByOrigin)
		parsedPathRoute.GET("/paths/metadata", getPathsMetadata)
	}
}

func getPathsMetadata(c *gin.Context) {
	projectName := c.Param("projectName")
	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	project := projectDb.GetOneProjectByName(projectName)
	if project == nil {
		responser.Response404(c, "Not found project with this name")
		return
	}

	var pathDb database.IParsedPathDatabase = new(database.ParsedPathDatabase).Init()
	origins, numPaths, numScanned, numFound, scanning := pathDb.GetMetadataByProject(project)
	result := make(map[string]interface{})
	result["origins"] = origins
	result["numPaths"] = numPaths
	result["numScanned"] = numScanned
	result["numFound"] = numFound
	result["scanning"] = scanning
	responser.ResponseOk(c, result)
}

func getPathsByOrigin(c *gin.Context) {
	projectName := c.Param("projectName")
	origin := c.Query("origin")

	var pathDb database.IParsedPathDatabase = new(database.ParsedPathDatabase).Init()
	paths := pathDb.GetPathsByProjectAndOrigin(projectName, origin)
	responser.ResponseOk(c, paths)
}
