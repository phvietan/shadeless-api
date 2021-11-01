package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/database/schema"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
)

func BotFuzzerRoutes(route *gin.Engine) {
	BotFuzzerRoute := route.Group("/projects/:projectName")
	{
		BotFuzzerRoute.GET("/bot_fuzzer", getBotFuzzerByProject)
		BotFuzzerRoute.PUT("/bot_fuzzer", putBotFuzzerByProject)
		BotFuzzerRoute.PUT("/bot_fuzzer/run", switchBotFuzzerRunning)
	}
}

func getBotFuzzerByProject(c *gin.Context) {
	projectName := c.Param("projectName")

	var BotFuzzerDb database.IBotFuzzerDatabase = new(database.BotFuzzerDatabase).Init()
	BotFuzzer := BotFuzzerDb.GetBotFuzzerByProject(projectName)
	if BotFuzzer == nil {
		responser.Response404(c, "BotFuzzer not found")
		return
	}

	responser.ResponseOk(c, BotFuzzer)
}

func switchBotFuzzerRunning(c *gin.Context) {
	projectName := c.Param("projectName")

	var BotFuzzerDb database.IBotFuzzerDatabase = new(database.BotFuzzerDatabase).Init()
	BotFuzzer := BotFuzzerDb.GetBotFuzzerByProject(projectName)
	if BotFuzzer == nil {
		responser.Response404(c, "BotFuzzer not found")
		return
	}

	if err := BotFuzzerDb.SwitchRun(BotFuzzer); err != nil {
		responser.ResponseError(c, err)
		return
	}
	msg := "Bot path set to running"
	if BotFuzzer.Running == true {
		msg = "Bot path stop running"
	}
	responser.ResponseOk(c, msg)
}

func putBotFuzzerByProject(c *gin.Context) {
	projectName := c.Param("projectName")

	var BotFuzzerDb database.IBotFuzzerDatabase = new(database.BotFuzzerDatabase).Init()
	BotFuzzer := BotFuzzerDb.GetBotFuzzerByProject(projectName)
	if BotFuzzer == nil {
		responser.Response404(c, "BotFuzzer not found")
		return
	}

	newBotFuzzer := schema.NewBotFuzzer(projectName)
	if err := c.BindJSON(newBotFuzzer); err != nil {
		responser.ResponseError(c, err)
		return
	}
	if err := BotFuzzerDb.PutBotFuzzerByProject(BotFuzzer.ID, newBotFuzzer); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Sucessfully update bot_path")
}
