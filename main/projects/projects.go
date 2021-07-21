package projects

import (
	"errors"
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
	ProjectPacketRoutes(route)
}

func getProjects(c *gin.Context) {
	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	projects := projectDb.GetProjects()
	responser.ResponseOk(c, projects)
}

func isProjectExist(name string) bool {
	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	check := projectDb.GetOneProjectByName(name)
	return check != nil
}

func postProjects(c *gin.Context) {
	project := database.NewProject()
	if err := c.BindJSON(project); err != nil {
		responser.ResponseError(c, err)
		return
	}

	if isProjectExist(project.Name) {
		responser.ResponseError(c, errors.New("Project with that name is already exist"))
		return
	}

	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	if err := projectDb.CreateProject(project); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully create project")
}

func putProjects(c *gin.Context) {
	newProject := database.NewProject()
	if err := c.BindJSON(newProject); err != nil {
		responser.ResponseError(c, err)
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		responser.ResponseError(c, err)
		return
	}

	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	dbProject := projectDb.GetOneProjectById(id)
	if dbProject.Name != newProject.Name {
		if isProjectExist(newProject.Name) {
			responser.ResponseError(c, errors.New("Project with that name is already exist"))
			return
		}
	}

	if err = projectDb.UpdateProject(id, newProject); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully update project")
}

type deleteOption struct {
	All bool `json:"all"`
}

func deleteProjects(c *gin.Context) {
	option := new(deleteOption)
	if err := c.BindJSON(option); err != nil {
		responser.ResponseError(c, err)
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		responser.ResponseError(c, err)
		return
	}

	var projectDb database.IProjectDatabase = new(database.ProjectDatabase).Init()
	project := projectDb.GetOneProjectById(id)
	if project == nil {
		responser.ResponseError(c, errors.New("Cannot delete project: project with this id is not found"))
		return
	}

	if err = projectDb.DeleteProject(id); err != nil {
		responser.ResponseError(c, err)
		return
	}
	if option.All == true {
		var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
		if err := packetDb.DeletePacketsByProjectName(project.Name); err != nil {
			responser.ResponseError(c, err)
			return
		}
	}
	responser.ResponseJson(c, 200, "Successfully delete project", "")
}
