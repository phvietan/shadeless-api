package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Routes(route *gin.Engine) {
	projects := route.Group("/projects")
	{
		projects.GET("/", getProjects)
		projects.POST("/", postProjects)
		projects.PUT("/:id", putProjects)
		projects.DELETE("/:id", deleteProjects)
	}
}

func getProjects(c *gin.Context) {
	projects := database.GetProjects()
	responser.ResponseOk(c, projects)
}

func postProjects(c *gin.Context) {
	project := database.NewProject()
	c.BindJSON(project)
	err := database.CreateProject(project)
	if err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully create project")
}

func putProjects(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		responser.ResponseError(c, err)
	}
	project := database.NewProject()
	c.BindJSON(project)
	err = database.UpdateProject(id, project)
	if err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully update project")
}

func deleteProjects(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		responser.ResponseError(c, err)
	}
	err = database.DeleteProject(id)
	if err != nil {
		responser.ResponseJson(c, 500, "", "Cannot delete project: "+err.Error())
		return
	}
	responser.ResponseJson(c, 200, "Successfully delete project", "")
}
