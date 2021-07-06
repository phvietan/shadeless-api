package projects

import (
	"net/http"
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
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
	responser.ResponseJson(c, http.StatusOK, projects, "")
}

func postProjects(c *gin.Context) {
	project := database.NewProject()
	c.BindJSON(project)

	err := database.CreateProject(project)
	if err != nil {
		responser.ResponseJson(c, 500, "", "Cannot create project: "+err.Error())
		return
	}

	responser.ResponseJson(c, 200, "Successfully create project", "")
}

func putProjects(c *gin.Context) {

}

func deleteProjects(c *gin.Context) {

}
