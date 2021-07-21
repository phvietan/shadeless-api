package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/finder"
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
		projects.GET("/packets", getPacketsByOrigin)
		projects.GET("/numberPackets", getNumPacketsByOrigin)
	}
}

func getProjectMetadata(c *gin.Context) {
	projectName := c.Param("projectName")

	var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	origins := packetDb.GetOriginsByProjectName(projectName)
	parameters := packetDb.GetParametersByProjectName(projectName)
	reflectedParameters := packetDb.GetReflectedParametersByProjectName(projectName)

	metaData := NewMetaData(origins, parameters, reflectedParameters)
	responser.ResponseOk(c, metaData)
}

func getNumPacketsByOrigin(c *gin.Context) {
	projectName := c.Param("projectName")
	origin := c.Query("origin")

	var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	numPackets := packetDb.GetNumPacketsByOrigin(projectName, origin)
	responser.ResponseOk(c, numPackets)
}

func getPacketsByOrigin(c *gin.Context) {
	projectName := c.Param("projectName")
	origin := c.Query("origin")

	options := new(finder.FinderOptions)
	err := c.Bind(options)
	if err != nil {
		responser.ResponseError(c, err)
		return
	}

	var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	packets := packetDb.GetPacketsByOriginAndProject(projectName, origin, options)
	responser.ResponseOk(c, packets)
}
