package database

import (
	"fmt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type PacketDatabase struct {
	Database
}

func (this *PacketDatabase) Init() *PacketDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&Packet{})
	return this
}

func getUniqueOriginsFromPackets(arr []Packet) []string {
	checkMap := make(map[string]bool)
	result := make([]string, 0)
	for _, packet := range arr {
		if !checkMap[packet.Origin] {
			checkMap[packet.Origin] = true
			result = append(result, packet.Origin)
		}
	}
	return result
}
func getUniqueParametersFromPackets(arr []Packet) []string {
	checkMap := make(map[string]bool)
	result := make([]string, 0)
	for _, packet := range arr {
		for _, param := range packet.Parameters {
			if !checkMap[param] {
				checkMap[param] = true
				result = append(result, param)
			}
		}
	}
	return result
}
func getReflectedParametersFromPackets(arr []Packet) map[string]string {
	result := make(map[string]string)
	for _, packet := range arr {
		for key, val := range packet.ReflectedParameters {
			result[key] = val
		}
	}
	return result
}

func parseFilterOptionsFromProject(project *Project) bson.M {
	blacklistExact := make([]string, 0)
	blacklistRegex := ""
	for _, bl := range project.Blacklist {
		if bl.Type == BlacklistValue {
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

func (this *PacketDatabase) GetPacketsAsTimeTravel(projectName string, packetPrefix string, packetIndex int, number int) []Packet {
	pipeline := []bson.M{
		bson.M{"$match": bson.M{
			"requestPacketPrefix": packetPrefix,
			"project":             projectName,
			"requestPacketIndex":  bson.M{"$gte": packetIndex, "$lt": packetIndex + number},
		}},
		bson.M{"$sort": bson.M{"requestPacketIndex": 1}},
		bson.M{"$project": bson.M{
			"requestPacketId":     1,
			"requestPacketPrefix": 1,
			"requestPacketIndex":  1,
			"requestBodyHash":     1,
			"responseBodyHash":    1,
			"requestHeaders":      1,
			"responseHeaders":     1,
			"codename":            1,
			"reflectedParameters": 1,
			"created_at":          1,
			"origin":              1,
			"path":                1,
		}},
	}
	cursor, err := this.db.Aggregate(this.ctx, pipeline)
	if err != nil {
		fmt.Println("Error in GetPacketsAsTimeTravel1: ", err)
		return []Packet{}
	}

	results := []Packet{}
	if err := cursor.All(this.ctx, &results); err != nil {
		fmt.Println("Error in GetPacketsAsTimeTravel2: ", err)
		return []Packet{}
	}
	return results
}

func (this *PacketDatabase) GetPacketByPacketId(projectName string, packetId string) *Packet {
	result := &Packet{}
	if err := this.db.FirstWithCtx(
		this.ctx,
		bson.M{"project": projectName, "requestPacketId": packetId},
		result,
	); err != nil {
		fmt.Errorf("%v", err)
		return nil
	}
	return result
}
