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

type IParsedPacketDatabase interface {
	IDatabase

	Init() *ParsedPacketDatabase
	Upsert(packet *schema.ParsedPacket) error
	GetMetadataByProject(project *schema.Project) ([]string, []string, map[string]string)
	GetPacketsByOriginAndProject(projectName string, origin string) []schema.ParsedPacket
	GetParsedByRawPackets(project string, packets []schema.Packet) []schema.ParsedPacket

	UpdateParsedPacketScore(id primitive.ObjectID, newScore float64) error
	GetFuzzingPacketsApi(project *schema.Project) ([]schema.ParsedPacket, []schema.ParsedPacket, []schema.ParsedPacket)
	GetFuzzingPacketsStatic(project *schema.Project) []schema.ParsedPacket

	GetOneById(id primitive.ObjectID) (*schema.ParsedPacket, error)
	ResetStatus(id primitive.ObjectID) error
}

type ParsedPacketDatabase struct {
	Database
}

func (this *ParsedPacketDatabase) Init() *ParsedPacketDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&schema.ParsedPacket{})
	mod := mongo.IndexModel{
		Keys: bson.D{
			{"project", 1},
			{"hash", 1},
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

func (this *ParsedPacketDatabase) GetOneById(id primitive.ObjectID) (*schema.ParsedPacket, error) {
	result := new(schema.ParsedPacket)
	if err := this.db.FindByID(id, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *ParsedPacketDatabase) ResetStatus(id primitive.ObjectID) error {
	_, err := this.db.UpdateByID(this.ctx, id, bson.M{
		"$set": bson.M{
			"status": schema.FuzzStatusTodo,
			"logDir": "",
		},
	})
	return err
}

func (this *ParsedPacketDatabase) UpdateParsedPacketScore(id primitive.ObjectID, newScore float64) error {
	_, err := this.db.UpdateByID(this.ctx, id, bson.M{
		"$set": bson.M{
			"staticScore": newScore,
		},
	})
	return err
}

func (this *ParsedPacketDatabase) GetFuzzingPacketsApi(project *schema.Project) ([]schema.ParsedPacket, []schema.ParsedPacket, []schema.ParsedPacket) {
	if project == nil {
		fmt.Println("Error project is nil")
		return []schema.ParsedPacket{}, []schema.ParsedPacket{}, []schema.ParsedPacket{}
	}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})
	filter := ParseFilterOptionsFromProject(project)
	filter["staticScore"] = bson.M{"$lte": 50}

	filter["status"] = schema.FuzzStatusDone
	done := []schema.ParsedPacket{}
	if err := this.db.SimpleFind(&done, filter, findOptions); err != nil {
		fmt.Println(err)
		return []schema.ParsedPacket{}, []schema.ParsedPacket{}, []schema.ParsedPacket{}
	}

	filter["status"] = schema.FuzzStatusScanning
	scanning := []schema.ParsedPacket{}
	if err := this.db.SimpleFind(&scanning, filter, findOptions); err != nil {
		fmt.Println(err)
		return []schema.ParsedPacket{}, []schema.ParsedPacket{}, []schema.ParsedPacket{}
	}

	filter["status"] = schema.FuzzStatusTodo
	todo := []schema.ParsedPacket{}
	if err := this.db.SimpleFind(&todo, filter, findOptions); err != nil {
		fmt.Println(err)
		return []schema.ParsedPacket{}, []schema.ParsedPacket{}, []schema.ParsedPacket{}
	}
	return done, scanning, todo
}

func (this *ParsedPacketDatabase) GetFuzzingPacketsStatic(project *schema.Project) []schema.ParsedPacket {
	if project == nil {
		fmt.Println("Error project is nil")
		return []schema.ParsedPacket{}
	}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})
	filter := ParseFilterOptionsFromProject(project)
	filter["staticScore"] = bson.M{"$gt": 50}

	results := []schema.ParsedPacket{}
	if err := this.db.SimpleFind(&results, filter, findOptions); err != nil {
		fmt.Println(err)
		return []schema.ParsedPacket{}
	}
	return results
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
	// If found, then update. However, other properties like status, static_score should not be updated
	packet.Status = result.Status
	packet.StaticScore = result.StaticScore
	_, err := this.db.UpdateByID(this.ctx, result.ID, bson.M{
		"$set": packet,
	})
	return err
}

func ParseFilterOptionsFromProject(project *schema.Project) bson.M {
	blacklistExact := make([]string, 0)
	blacklistRegex := ""
	for _, bl := range project.Blacklist {
		if bl.Type == schema.BlacklistValue {
			blacklistExact = append(blacklistExact, bl.Value)
		} else {
			blacklistRegex = bl.Value
		}
	}
	filter := bson.M{}
	if blacklistRegex != "" {
		filter = bson.M{
			"project": project.Name,
			"origin": bson.M{
				"$nin":   blacklistExact,
				"$not":   bson.M{"$regex": blacklistRegex},
				"$regex": project.Whitelist,
			},
		}
	} else {
		filter = bson.M{
			"project": project.Name,
			"origin": bson.M{
				"$nin":   blacklistExact,
				"$regex": project.Whitelist,
			},
		}
	}
	return filter
}

func (this *ParsedPacketDatabase) GetMetadataByProject(project *schema.Project) ([]string, []string, map[string]string) {
	if project == nil {
		return []string{}, []string{}, make(map[string]string)
	}
	filter := ParseFilterOptionsFromProject(project)
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

func (this *ParsedPacketDatabase) GetPacketsByOriginAndProject(projectName string, origin string) []schema.ParsedPacket {
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"origin": origin, "project": projectName}},
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

func (this *ParsedPacketDatabase) GetParsedByRawPackets(project string, packets []schema.Packet) []schema.ParsedPacket {
	packetsHash := []string{}
	for _, p := range packets {
		packetsHash = append(packetsHash, schema.CalculatePacketHash(
			p.Method,
			p.ResponseStatus,
			p.Origin,
			p.Path,
			p.Parameters,
		))
	}
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"hash": bson.M{
					"$in": packetsHash,
				},
				"project": project,
			},
		},
	}
	cursor, err := this.db.Aggregate(this.ctx, pipeline)
	if err != nil {
		fmt.Println("Error in GetParsedByRawPackets: ", err)
		return []schema.ParsedPacket{}
	}

	allParsedPacketInDB := []schema.ParsedPacket{}
	if err := cursor.All(this.ctx, &allParsedPacketInDB); err != nil {
		fmt.Println("Error in GetParsedByRawPackets: ", err)
		return []schema.ParsedPacket{}
	}

	result := []schema.ParsedPacket{}
	for idx, h := range packetsHash {
		found := false
		for _, parsedPacket := range allParsedPacketInDB {
			if h == parsedPacket.Hash {
				curP := parsedPacket
				curP.RequestPacketId = packets[idx].RequestPacketId
				curP.RequestPacketIndex = packets[idx].RequestPacketIndex
				curP.RequestPacketPrefix = packets[idx].RequestPacketPrefix
				result = append(result, curP)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Not found time travel: ", packets[idx].RequestPacketId)
		}
	}
	return result
}
