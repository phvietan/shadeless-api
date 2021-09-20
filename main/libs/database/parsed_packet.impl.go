package database

import (
	"errors"
	"fmt"
	"shadeless-api/main/libs"
	"shadeless-api/main/libs/database/schema"
	"shadeless-api/main/libs/finder"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type ParsedPacketDatabase struct {
	Database
}

func (this *ParsedPacketDatabase) Init() *ParsedPacketDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&schema.ParsedPacket{})
	return this
}

func (this *ParsedPacketDatabase) Upsert(packet *schema.ParsedPacket) error {
	if packet == nil {
		return errors.New("Parsed packet object is nil")
	}
	result := &schema.ParsedPacket{}
	if err := this.db.FirstWithCtx(
		this.ctx,
		bson.M{"hash": packet.Hash, "project": packet.Project},
		result,
	); err != nil {
		// If not found, then insert
		fmt.Println("Not found parse packet, creating one")
		return this.Insert(packet)
	}
	// If found, then update, the fuzzed property should not be update, lol
	packet.Fuzzed = result.Fuzzed
	_, err := this.db.UpdateByID(this.ctx, result.ID, bson.M{
		"$set": packet,
	})
	return err
}

func (this *ParsedPacketDatabase) GetMetadataByProject(project *schema.Project) ([]string, []string, map[string]string) {
	if project == nil {
		return []string{}, []string{}, make(map[string]string)
	}
	filter := parseFilterOptionsFromProject(project)
	resultOrigins, err := this.db.Distinct(this.ctx, "origin", filter)
	if err != nil {
		fmt.Println("Cannot get origins in metadata by project: ", err)
		return []string{}, []string{}, make(map[string]string)
	}
	origins := libs.ArrayInterfaceToArrayString(resultOrigins)

	resultParameters, err := this.db.Distinct(this.ctx, "parameters", filter)
	if err != nil {
		fmt.Println("Cannot get parameters in metadata by project: ", err)
		return []string{}, []string{}, make(map[string]string)
	}
	parameters := libs.ArrayInterfaceToArrayString(resultParameters)

	resultReflectedParameters, err := this.db.Distinct(this.ctx, "reflectedParameters", filter)
	if err != nil {
		fmt.Println("Cannot get reflected parameters in metadata by project: ", err)
		return []string{}, []string{}, make(map[string]string)
	}
	reflectedParameters := libs.ArrayInterfaceToMapString(resultReflectedParameters)
	return origins, parameters, reflectedParameters
}

func (this *ParsedPacketDatabase) GetPacketsByOriginAndProject(projectName string, origin string, options *finder.FinderOptions) []schema.ParsedPacket {
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"origin": origin, "project": projectName}},
		bson.M{"$group": bson.M{"_id": "$hash", "doc": bson.M{"$last": "$$ROOT"}}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$doc"}},
		bson.M{"$skip": options.Skip},
		bson.M{"$limit": options.Limit},
	}

	cursor, err := this.db.Aggregate(this.ctx, pipeline)
	if err != nil {
		fmt.Println("Error GetPacketsByOriginAndProject1", err)
		return []schema.ParsedPacket{}
	}

	results := []schema.ParsedPacket{}
	if err := cursor.All(this.ctx, &results); err != nil {
		fmt.Println("Error GetPacketsByOriginAndProject2", err)
		return []schema.ParsedPacket{}
	}
	return results
}
