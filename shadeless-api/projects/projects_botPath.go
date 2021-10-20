package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/database/schema"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
)

func BotPathRoutes(route *gin.Engine) {
	botPathRoute := route.Group("/projects/:projectName")
	{
		botPathRoute.GET("/bot_path", getBotPathByProject)
		botPathRoute.PUT("/bot_path", putBotPathByProject)
		botPathRoute.PUT("/bot_path/run", switchBotPathRunning)
	}
}

func getBotPathByProject(c *gin.Context) {
	projectName := c.Param("projectName")

	var botPathDb database.IBotPathDatabase = new(database.BotPathDatabase).Init()
	botPath := botPathDb.GetBotPathByProject(projectName)
	if botPath == nil {
		responser.Response404(c, "BotPath not found")
		return
	}

	responser.ResponseOk(c, botPath)
}

func switchBotPathRunning(c *gin.Context) {
	projectName := c.Param("projectName")

	var botPathDb database.IBotPathDatabase = new(database.BotPathDatabase).Init()
	botPath := botPathDb.GetBotPathByProject(projectName)
	if botPath == nil {
		responser.Response404(c, "BotPath not found")
		return
	}

	if err := botPathDb.SwitchRun(botPath); err != nil {
		responser.ResponseError(c, err)
		return
	}
	msg := "Bot path set to running"
	if botPath.Running == true {
		msg = "Bot path stop running"
	}
	responser.ResponseOk(c, msg)
}

func putBotPathByProject(c *gin.Context) {
	projectName := c.Param("projectName")

	var botPathDb database.IBotPathDatabase = new(database.BotPathDatabase).Init()
	botPath := botPathDb.GetBotPathByProject(projectName)
	if botPath == nil {
		responser.Response404(c, "BotPath not found")
		return
	}

	newBotPath := schema.NewBotPath(projectName)
	if err := c.BindJSON(newBotPath); err != nil {
		responser.ResponseError(c, err)
		return
	}
	if err := botPathDb.PutBotPathByProject(botPath.ID, newBotPath); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Sucessfully update bot_path")
}
