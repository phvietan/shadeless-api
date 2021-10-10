package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(route *gin.Engine) {
	users := route.Group("/projects/:projectName")
	{
		users.GET("/users", getUsersInProject)
	}
}

func getUsersInProject(c *gin.Context) {
	var userDb database.IUserDatabase = new(database.UserDatabase).Init()
	users := userDb.GetUsers(c.Param("projectName"))
	responser.ResponseOk(c, users)
}
