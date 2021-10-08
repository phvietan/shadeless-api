package database

import (
	"errors"
	"fmt"
	"shadeless-api/main/libs"
	"shadeless-api/main/libs/database/schema"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type IParsedPacketDatabase interface {
	IDatabase

	Init() *ParsedPacketDatabase
	Upsert(packet *schema.ParsedPacket) error
	GetMetadataByProject(project *schema.Project) ([]string, []string, map[string]string)
	GetPacketsByOriginAndProject(projectName string, origin string) []schema.ParsedPacket
	GetParsedByRawPackets(project string, packets []schema.Packet) []schema.ParsedPacket
}

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
	// If found, then update, other properties like fuzzed, static_score should not be updated
	packet.Fuzzed = result.Fuzzed
	_, err := this.db.UpdateByID(this.ctx, result.ID, bson.M{
		"$set": packet,
	})
	return err
}

func parseFilterOptionsFromProject(project *schema.Project) bson.M {
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
				"$nin": blacklistExact,
				"$not": bson.M{"$regex": blacklistRegex},
			},
		}
	} else {
		filter = bson.M{
			"project": project.Name,
			"origin": bson.M{
				"$nin": blacklistExact,
			},
		}
	}
	return filter
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
				result = append(result, parsedPacket)
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
