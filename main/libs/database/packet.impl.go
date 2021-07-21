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
			"codeName":            1,
			"reflectedParameters": 1,
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
