package database

import (
	"fmt"
	"shadeless-api/main/libs/database/schema"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoteDatabase struct {
	Database
}

func (this *NoteDatabase) Init() *NoteDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&schema.Note{})
	return this
}

func (this *NoteDatabase) GetNotes(project string) []schema.Note {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	results := []schema.Note{}
	if err := this.db.SimpleFind(&results, bson.M{
		"project": project,
	}, findOptions); err != nil {
		fmt.Println(err)
		return []schema.Note{}
	}
	return results
}

func (this *NoteDatabase) UpdateOne(id primitive.ObjectID, newNote *schema.Note) error {
	updated := bson.M{
		"description": newNote.Description,
		"tags":        newNote.Tags,
		"codeName":    newNote.CodeName,
	}
	if _, err := this.db.UpdateByID(this.ctx, id, bson.D{{"$set", updated}}); err != nil {
		return err
	}
	return nil
}

func (this *NoteDatabase) GetNotesByPackets(project string, packets []schema.Packet) []*schema.Note {
	packetsId := []string{}
	for _, p := range packets {
		packetsId = append(packetsId, p.RequestPacketId)
	}
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"requestPacketId": bson.M{
					"$in": packetsId,
				},
				"project": project,
			},
		},
	}
	cursor, err := this.db.Aggregate(this.ctx, pipeline)
	if err != nil {
		fmt.Println("Error in GetNotesByPackets: ", err)
		ret := []*schema.Note{}
		for i := 0; i < len(packets); i++ {
			ret = append(ret, nil)
		}
		return ret
	}

	notes := []schema.Note{}
	if err := cursor.All(this.ctx, &notes); err != nil {
		fmt.Println("Error in GetNotesByPackets: ", err)
		ret := []*schema.Note{}
		for i := 0; i < len(packets); i++ {
			ret = append(ret, nil)
		}
		return ret
	}

	results := []*schema.Note{}
	for _, p := range packets {
		var found *schema.Note = nil
		for _, note := range notes {
			if note.RequestPacketId == p.RequestPacketId {
				found = &note
			}
		}
		results = append(results, found)
	}
	return results
}
