package projects

import (
	"errors"
	"fmt"
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/finder"
	"shadeless-api/main/libs/responser"
	"strconv"
	"strings"

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
		projects.GET("", getProjectByName)
		projects.GET("/metadata", getProjectMetadata)
		projects.GET("/packets", getPacketsByOrigin)
		projects.GET("/numberPackets", getNumPacketsByOrigin)
		projects.GET("/timeTravel", getTimeTravel)
	}
}

func getProjectByName(c *gin.Context) {
	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	project := projectDb.GetOneProjectByName(c.Param("projectName"))
	responser.ResponseOk(c, project)
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
	if err := c.BindQuery(options); err != nil {
		responser.ResponseError(c, err)
		return
	}

	var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	packets := packetDb.GetPacketsByOriginAndProject(projectName, origin, options)
	responser.ResponseOk(c, packets)
}

type timeTravelOptions struct {
	RequestPacketId string `form:"requestPacketId"`
	Number          int    `form:"number"`
}

func getTimeTravel(c *gin.Context) {
	fmt.Println("???")
	options := new(timeTravelOptions)
	if err := c.Bind(options); err != nil {
		responser.ResponseError(c, err)
		return
	}
	projectName := c.Param("projectName")
	arr := strings.Split(options.RequestPacketId, ".")
	if len(arr) != 2 {
		fmt.Println("1: ", options.RequestPacketId)
		responser.ResponseError(c, errors.New("Wrong requestPacketId format"))
		return
	}
	var packetIndex int
	var err error
	if packetIndex, err = strconv.Atoi(arr[1]); err != nil {
		fmt.Println("2: ", options.RequestPacketId)
		responser.ResponseError(c, errors.New("Wrong requestPacketId format"))
		return
	}
	packetPrefix := arr[0]

	var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	packets := packetDb.GetPacketsAsTimeTravel(projectName, packetPrefix, packetIndex, options.Number)
	responser.ResponseOk(c, packets)
}
