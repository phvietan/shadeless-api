package projects

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/database/schema"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NotesRoutes(route *gin.Engine) {
	users := route.Group("/projects/:projectName")
	{
		users.GET("/notes", getNotesInProject)
		users.POST("/notes", createNewNote)
		users.PUT("/notes/:noteId", updateNote)
	}
}

func getNotesInProject(c *gin.Context) {
	var noteDb database.INoteDatabase = new(database.NoteDatabase).Init()
	users := noteDb.GetNotes(c.Param("projectName"))
	responser.ResponseOk(c, users)
}

func createNewNote(c *gin.Context) {
	note := schema.NewNote()
	if err := c.BindJSON(note); err != nil {
		responser.ResponseError(c, err)
		return
	}
	note.Project = c.Param("projectName")
	var projectDb database.INoteDatabase = new(database.NoteDatabase).Init()
	if err := projectDb.Insert(note); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully create note")
}

func updateNote(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("noteId"))
	if err != nil {
		responser.ResponseError(c, err)
		return
	}
	newNote := schema.NewNote()
	if err := c.BindJSON(newNote); err != nil {
		responser.ResponseError(c, err)
		return
	}
	var noteDb database.INoteDatabase = new(database.NoteDatabase).Init()
	if err := noteDb.UpdateOne(id, newNote); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully update note")
}
