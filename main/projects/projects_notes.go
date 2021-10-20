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
		users.DELETE("/notes/:noteId", deleteNote)
	}
}

func getNotesInProject(c *gin.Context) {
	var noteDb database.INoteDatabase = new(database.NoteDatabase).Init()
	projectName := c.Param("projectName")
	notes := noteDb.GetNotes(projectName)

	var packetsDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	packets := packetsDb.GetPacketsByNotes(projectName, notes)

	result := make(map[string]interface{})
	result["notes"] = notes
	result["packets"] = packets

	responser.ResponseOk(c, result)
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

func deleteNote(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("noteId"))
	if err != nil {
		responser.ResponseError(c, err)
		return
	}
	var noteDb database.INoteDatabase = new(database.NoteDatabase).Init()
	if err := noteDb.DeleteById(id); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully delete note")
}
