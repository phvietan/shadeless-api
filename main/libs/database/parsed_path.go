package database

import (
	"errors"
	"fmt"
	"os"
	"shadeless-api/main/libs/database/schema"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IParsedPathDatabase interface {
	IDatabase

	Init() *ParsedPathDatabase
	GetPathsByProject(project string) []schema.ParsedPath
	Upsert(parsedPath *schema.ParsedPath) error
}

type ParsedPathDatabase struct {
	Database
}

func (this *ParsedPathDatabase) Init() *ParsedPathDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&schema.ParsedPath{})
	mod := mongo.IndexModel{
		Keys: bson.D{
			{"project", 1},
			{"origin", 1},
			{"path", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := this.db.Indexes().CreateOne(this.ctx, mod)
	if err != nil {
		fmt.Println("Error when creating index, ", err)
		os.Exit(0)
	}
	return this
}

func (this *ParsedPathDatabase) GetPathsByProject(project string) []schema.ParsedPath {
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"project": project}},
	}

	cursor, err := this.db.Aggregate(this.ctx, pipeline)
	if err != nil {
		fmt.Println("Error GetPacketsByOriginAndProject1", err)
		return []schema.ParsedPath{}
	}

	results := []schema.ParsedPath{}
	if err := cursor.All(this.ctx, &results); err != nil {
		fmt.Println("Error GetPacketsByOriginAndProject2", err)
		return []schema.ParsedPath{}
	}
	return results
}

func (this *ParsedPathDatabase) Upsert(parsedPath *schema.ParsedPath) error {
	if parsedPath == nil {
		fmt.Println("Error: parsedPath is nil, cannot upsert")
		return errors.New("parsedPath is nil, cannot upsert")
	}
	result := &schema.ParsedPacket{}
	if err := this.db.FirstWithCtx(
		this.ctx,
		bson.M{"origin": parsedPath.Origin, "project": parsedPath.Project},
		result,
	); err != nil {
		// If not found, then insert
		fmt.Println("Not found parsedPath, creating one")
		return this.Insert(parsedPath)
	}
	// If found, then update requestPacketId
	_, err := this.db.UpdateByID(this.ctx, result.ID, bson.M{
		"$set": bson.M{"requestPacketId": parsedPath.RequestPacketId},
	})
	return err
}
