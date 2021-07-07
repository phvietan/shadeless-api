package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
)

type metaData struct {
	Origins             []string `json:"origins"`
	Parameters          []string `json:"parameters"`
	ReflectedParameters []string `json:"reflectedParameters"`
}

func NewMetaData(origins []string, parameters []string, reflectedParameters []string) *metaData {
	return &metaData{
		Origins:             origins,
		Parameters:          parameters,
		ReflectedParameters: reflectedParameters,
	}
}

func ProjectPacketRoutes(route *gin.Engine) {
	projects := route.Group("/projects/:projectName")
	{
		projects.GET("/metadata", getProjectMetadata)
	}
}

func getProjectMetadata(c *gin.Context) {
	projectName := c.Param("projectName")
	origins := database.GetOrigins(projectName)
	parameters := database.GetParameters(projectName)
	reflectedParameters := database.GetReflectedParameters(projectName)

	metaData := NewMetaData(origins, parameters, reflectedParameters)
	responser.ResponseOk(c, metaData)
}

func getParameters(c *gin.Context) {

}
