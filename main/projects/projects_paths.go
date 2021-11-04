package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/database/schema"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PathsRoutes(route *gin.Engine) {
	parsedPathRoute := route.Group("/projects/:projectName")
	{
		parsedPathRoute.GET("/paths", getParsedPaths)
		parsedPathRoute.PUT("/paths/:id/status", putPathStatus)
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

type pathStatus struct {
	Status string `json:"status"`
}

func putPathStatus(c *gin.Context) {
	projectName := c.Param("projectName")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		responser.ResponseError(c, err)
		return
	}

	pathStatus := new(pathStatus)
	if err := c.BindJSON(pathStatus); err != nil {
		responser.ResponseError(c, err)
		return
	}

	var pathDb database.IParsedPathDatabase = new(database.ParsedPathDatabase).Init()
	if err := pathDb.UpdateStatus(projectName, id, pathStatus.Status); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully update path status to "+pathStatus.Status)
}

func getParsedPaths(c *gin.Context) {
	projectName := c.Param("projectName")

	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	project := projectDb.GetOneProjectByName(projectName)
	if project == nil {
		responser.Response404(c, "No such project")
	}

	var pathDb database.IParsedPathDatabase = new(database.ParsedPathDatabase).Init()
	filter := database.ParseFilterOptionsFromProject(project)
	filter["status"] = schema.FuzzStatusDone
	donePaths := pathDb.GetParsedPaths(filter)

	filter["status"] = schema.FuzzStatusTodo
	todoPaths := pathDb.GetParsedPaths(filter)

	filter["status"] = schema.FuzzStatusScanning
	scanningPaths := pathDb.GetParsedPaths(filter)

	filter["status"] = schema.FuzzStatusRemoved
	removedPaths := pathDb.GetParsedPaths(filter)

	result := make(map[string]interface{})
	result["done"] = donePaths
	result["todo"] = todoPaths
	result["scanning"] = scanningPaths
	result["removed"] = removedPaths

	responser.ResponseOk(c, result)
}
