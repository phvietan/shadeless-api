package projects

import (
	"errors"
	"fmt"
	"os"
	"path"
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/responser"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PacketsRoutes(route *gin.Engine) {
	projects := route.Group("/projects/:projectName")
	{
		projects.GET("", getProjectByName)
		projects.GET("/metadata", getProjectMetadata)
		projects.GET("/packets", getPacketsByOrigin)
		projects.GET("/timeTravel", getTimeTravel)
		projects.GET("/fuzzing_packets/api", getFuzzingPacketsApi)
		projects.GET("/fuzzing_packets/static", getFuzzingPacketsStatic)
		projects.PUT("/fuzzing_packets/:id/score", putFuzzingPacketScore)
		projects.PUT("/fuzzing_packets/:id/reset", putFuzzingPacketStatus)
	}
}

type putFuzzingPacketScoreObj struct {
	Score float64 `json:"score"`
}

func putFuzzingPacketScore(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		responser.ResponseError(c, err)
		return
	}
	newScore := new(putFuzzingPacketScoreObj)
	if err := c.BindJSON(newScore); err != nil {
		fmt.Println("Cannot bindJSON newScore object: ", err)
		responser.ResponseError(c, err)
		return
	}
	var parsedPacketDb database.IParsedPacketDatabase = new(database.ParsedPacketDatabase).Init()
	if err := parsedPacketDb.UpdateParsedPacketScore(id, newScore.Score); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully update score")
}

func putFuzzingPacketStatus(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		responser.ResponseError(c, err)
		return
	}
	var parsedPacketDb database.IParsedPacketDatabase = new(database.ParsedPacketDatabase).Init()
	parsedPacket, err := parsedPacketDb.GetOneById(id)
	if err != nil {
		responser.ResponseError(c, err)
		return
	}
	pathDelete := path.Join("../shadeless-bot/", parsedPacket.LogDir, "..")
	if err := os.RemoveAll(pathDelete); err != nil {
		responser.ResponseError(c, err)
		return
	}

	if err := parsedPacketDb.ResetStatus(id); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully reset status for this api")
}

func getProjectByName(c *gin.Context) {
	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	project := projectDb.GetOneProjectByName(c.Param("projectName"))
	if project == nil {
		responser.Response404(c, "Not found project with that name")
		return
	}
	responser.ResponseOk(c, project)
}

func getProjectMetadata(c *gin.Context) {
	projectName := c.Param("projectName")
	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	project := projectDb.GetOneProjectByName(projectName)
	if project == nil {
		responser.Response404(c, "Not found project with this name")
		return
	}
	var parsedPacketDb database.IParsedPacketDatabase = new(database.ParsedPacketDatabase).Init()
	origins, parameters, reflectedParameters := parsedPacketDb.GetMetadataByProject(project)

	result := make(map[string]interface{})
	result["origins"] = origins
	result["parameters"] = parameters
	result["reflectedParameters"] = reflectedParameters
	responser.ResponseOk(c, result)
}

func getFuzzingPacketsApi(c *gin.Context) {
	projectName := c.Param("projectName")
	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	project := projectDb.GetOneProjectByName(projectName)
	if project == nil {
		responser.Response404(c, "No such project")
	}

	var packetDb database.IParsedPacketDatabase = new(database.ParsedPacketDatabase).Init()
	done, scanning, todo := packetDb.GetFuzzingPacketsApi(project)

	result := make(map[string]interface{})
	result["done"] = done
	result["scanning"] = scanning
	result["todo"] = todo
	responser.ResponseOk(c, result)
}

func getFuzzingPacketsStatic(c *gin.Context) {
	projectName := c.Param("projectName")
	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	project := projectDb.GetOneProjectByName(projectName)
	if project == nil {
		responser.Response404(c, "No such project")
	}

	var packetDb database.IParsedPacketDatabase = new(database.ParsedPacketDatabase).Init()
	packets := packetDb.GetFuzzingPacketsStatic(project)
	responser.ResponseOk(c, packets)
}

func getPacketsByOrigin(c *gin.Context) {
	projectName := c.Param("projectName")
	origin := c.Query("origin")

	var packetDb database.IParsedPacketDatabase = new(database.ParsedPacketDatabase).Init()
	packets := packetDb.GetPacketsByOriginAndProject(projectName, origin)
	responser.ResponseOk(c, packets)
}

type timeTravelOptions struct {
	RequestPacketId string `form:"requestPacketId"`
	Number          int    `form:"number"`
}

func getTimeTravel(c *gin.Context) {
	options := new(timeTravelOptions)
	if err := c.Bind(options); err != nil {
		responser.ResponseError(c, err)
		return
	}
	projectName := c.Param("projectName")
	arr := strings.Split(options.RequestPacketId, ".")
	if len(arr) != 2 {
		responser.ResponseError(c, errors.New("Wrong requestPacketId format"))
		return
	}
	var packetIndex int
	var err error
	if packetIndex, err = strconv.Atoi(arr[1]); err != nil {
		responser.ResponseError(c, errors.New("Wrong requestPacketId format"))
		return
	}
	packetPrefix := arr[0]

	var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	packets := packetDb.GetPacketsAsTimeTravel(projectName, packetPrefix, packetIndex, options.Number)

	var parsedPacketDb database.IParsedPacketDatabase = new(database.ParsedPacketDatabase).Init()
	parsedPacket := parsedPacketDb.GetParsedByRawPackets(projectName, packets)

	if len(packets) != len(parsedPacket) {
		fmt.Println("Soemthing is wrong in time travel")
	}

	var noteDb database.INoteDatabase = new(database.NoteDatabase).Init()
	notes := noteDb.GetNotesByPackets(projectName, packets)

	result := make(map[string]interface{})
	result["packets"] = parsedPacket
	result["notes"] = notes
	responser.ResponseOk(c, result)
}
