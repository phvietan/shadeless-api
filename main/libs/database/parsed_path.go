package database

import (
	"errors"
	"fmt"
	"os"
	"shadeless-api/main/libs"
	"shadeless-api/main/libs/database/schema"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IParsedPathDatabase interface {
	IDatabase

	Init() *ParsedPathDatabase
	GetParsedPaths(filter bson.M) []schema.ParsedPath
	GetMetadataByProject(project *schema.Project) ([]string, int, int, int, string)
	Upsert(parsedPath *schema.ParsedPath) error

	UpdateStatus(projectName string, id primitive.ObjectID, newStatus string) error
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

// Get origins, num paths, num scanned, num found, who is bot scanning
func (this *ParsedPathDatabase) GetMetadataByProject(project *schema.Project) ([]string, int, int, int, string) {
	if project == nil {
		fmt.Println("Error ParsedPath.GetMetadataByProject: project is nil")
		return []string{}, 0, 0, 0, ""
	}
	filter := ParseFilterOptionsFromProject(project)
	resultOrigins, err := this.db.Distinct(this.ctx, "origin", filter)
	if err != nil {
		fmt.Println("Error ParsedPath.GetMetadataByProject:", err)
		return []string{}, 0, 0, 0, ""
	}
	origins := libs.ArrayInterfaceToArrayString(resultOrigins)

	numPaths, err := this.db.CountDocuments(this.ctx, filter)
	if err != nil {
		fmt.Println("Error ParsedPath.GetMetadataByProject", err)
		return []string{}, 0, 0, 0, ""
	}
	filter["requestPacketId"] = ""
	numFound, err := this.db.CountDocuments(this.ctx, filter)
	if err != nil {
		fmt.Println("Error ParsedPath.GetMetadataByProject", err)
		return []string{}, 0, 0, 0, ""
	}
	delete(filter, "requestPacketId")

	filter["status"] = schema.FuzzStatusDone
	numScanned, err := this.db.CountDocuments(this.ctx, filter)
	if err != nil {
		fmt.Println("Error ParsedPath.GetMetadataByProject2", err)
		return []string{}, 0, 0, 0, ""
	}
	botScanning := &schema.ParsedPath{}
	filter["status"] = schema.FuzzStatusScanning
	this.db.FirstWithCtx(this.ctx, filter, botScanning)
	return origins, int(numPaths), int(numScanned), int(numFound), botScanning.Origin + botScanning.Path
}

func (this *ParsedPathDatabase) GetParsedPaths(filter bson.M) []schema.ParsedPath {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	results := []schema.ParsedPath{}
	if err := this.db.SimpleFind(&results, filter, findOptions); err != nil {
		fmt.Println("Error: GetParsedPaths", err)
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
		bson.M{"origin": parsedPath.Origin, "project": parsedPath.Project, "path": parsedPath.Path},
		result,
	); err != nil {
		// If not found, then insert
		fmt.Println("Not found parsedPath, creating one")
		return this.Insert(parsedPath)
	}
	// If found, then update requestPacketId
	_, err := this.db.UpdateByID(this.ctx, result.ID, bson.M{
		"$set": bson.M{
			"requestPacketId": parsedPath.RequestPacketId,
		},
	})
	return err
}

func (this *ParsedPathDatabase) UpdateStatus(projectName string, id primitive.ObjectID, newStatus string) error {
	if newStatus != schema.FuzzStatusRemoved && newStatus != schema.FuzzStatusTodo {
		return errors.New("ParsedPath status is wrong")
	}
	filter := bson.M{"project": projectName, "_id": id}
	_, err := this.db.UpdateOne(this.ctx, filter, bson.M{
		"$set": bson.M{"status": newStatus},
	})
	return err
}
