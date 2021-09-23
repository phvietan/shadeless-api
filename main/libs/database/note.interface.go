package database

import (
	"shadeless-api/main/libs/database/schema"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type INoteDatabase interface {
	IDatabase
	Init() *NoteDatabase
	GetNotes(project string) []schema.Note
	UpdateOne(id primitive.ObjectID, newNote *schema.Note) error
	GetNotesByPackets(project string, packets []schema.Packet) []*schema.Note
}
