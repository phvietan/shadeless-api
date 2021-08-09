package database

import (
	"errors"
	"fmt"
	"log"
	"shadeless-api/main/libs"
	"shadeless-api/main/libs/finder"

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

func (this *PacketDatabase) CreatePacket(packet *Packet) error {
	if packet == nil {
		return errors.New("Packet object is nil")
	}
	if err := this.db.Create(packet); err != nil {
		fmt.Errorf("%v", err)
		return err
	}
	return nil
}

func (this *PacketDatabase) getDistincColumnWithFilterOptions(name string, filterOptions bson.M) []string {
	results, err := this.db.Distinct(this.ctx, name, filterOptions)
	if err != nil {
		fmt.Errorf("%v", err)
		return []string{}
	}
	return libs.ArrayInterfaceToArrayString(results)
}

func (this *PacketDatabase) GetOriginsByProjectName(projectName string) []string {
	filterOptions := bson.M{
		"project": projectName,
	}
	return this.getDistincColumnWithFilterOptions("origin", filterOptions)
}

func (this *PacketDatabase) GetParametersByProjectName(projectName string) []string {
	filterOptions := bson.M{
		"project": projectName,
	}
	return this.getDistincColumnWithFilterOptions("parameters", filterOptions)
}

func (this *PacketDatabase) GetReflectedParametersByProjectName(projectName string) []string {
	filterOptions := bson.M{
		"project": projectName,
	}
	return this.getDistincColumnWithFilterOptions("reflectedParameters", filterOptions)
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
	fmt.Println(result)
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

func (this *PacketDatabase) GetMetadataByProject(project *Project) ([]string, []string, map[string]string) {
	if project == nil {
		return []string{}, []string{}, make(map[string]string)
	}
	filter := parseFilterOptionsFromProject(project)
	resultOrigins, err := this.db.Distinct(this.ctx, "origin", filter)
	if err != nil {
		fmt.Errorf("%v", err)
		return []string{}, []string{}, make(map[string]string)
	}
	origins := libs.ArrayInterfaceToArrayString(resultOrigins)

	resultParameters, err := this.db.Distinct(this.ctx, "parameters", filter)
	if err != nil {
		fmt.Errorf("%v", err)
		return []string{}, []string{}, make(map[string]string)
	}
	parameters := libs.ArrayInterfaceToArrayString(resultParameters)

	resultReflectedParameters, err := this.db.Distinct(this.ctx, "reflectedParameters", filter)
	if err != nil {
		fmt.Errorf("%v", err)
		return []string{}, []string{}, make(map[string]string)
	}
	reflectedParameters := libs.ArrayInterfaceToMapString(resultReflectedParameters)
	return origins, parameters, reflectedParameters
}

func (this *PacketDatabase) GetNumPacketsByOrigin(projectName string, origin string) int32 {
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"origin": origin, "project": projectName}},
		bson.M{"$group": bson.M{"_id": "$path"}},
		bson.M{"$count": "result"},
	}

	cursor, err := this.db.Aggregate(this.ctx, pipeline)
	if err != nil {
		fmt.Errorf("%v", err)
		return 0
	}
	var allDocs []bson.M
	if err = cursor.All(this.ctx, &allDocs); err != nil {
		log.Fatal(err)
		return 0
	}
	return allDocs[0]["result"].(int32) // Ugly hacks, TODO: refactor
}

func (this *PacketDatabase) GetPacketsByOriginAndProject(projectName string, origin string, options *finder.FinderOptions) []Packet {
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"origin": origin, "project": projectName}},
		bson.M{"$group": bson.M{"_id": "$path", "doc": bson.M{"$last": "$$ROOT"}}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$doc"}},
		bson.M{"$skip": options.Skip},
		bson.M{"$limit": options.Limit},
	}

	cursor, err := this.db.Aggregate(this.ctx, pipeline)
	if err != nil {
		fmt.Errorf("%v", err)
		return []Packet{}
	}

	results := []Packet{}
	if err := cursor.All(this.ctx, &results); err != nil {
		fmt.Errorf("%v", err)
		return []Packet{}
	}
	return results
}

func (this *PacketDatabase) DeletePacketsByProjectName(projectName string) error {
	if _, err := this.db.DeleteMany(this.ctx, bson.M{"project": projectName}); err != nil {
		return err
	}
	return nil
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
		fmt.Errorf("%v", err)
		return []Packet{}
	}

	results := []Packet{}
	if err := cursor.All(this.ctx, &results); err != nil {
		fmt.Errorf("%v", err)
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

func (this *PacketDatabase) UpdateProjectName(oldName string, newName string) error {
	_, err := this.db.UpdateMany(
		this.ctx,
		bson.M{"project": oldName},
		bson.D{{"$set", bson.M{"project": newName}}},
	)
	return err
}
