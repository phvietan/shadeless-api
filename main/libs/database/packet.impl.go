package database

import (
	"fmt"
	"shadeless-api/main/libs/database/schema"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type PacketDatabase struct {
	Database
}

func (this *PacketDatabase) Init() *PacketDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&schema.Packet{})
	return this
}

func (this *PacketDatabase) GetPacketsAsTimeTravel(projectName string, packetPrefix string, packetIndex int, number int) []schema.Packet {
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
		return []schema.Packet{}
	}

	results := []schema.Packet{}
	if err := cursor.All(this.ctx, &results); err != nil {
		fmt.Println("Error in GetPacketsAsTimeTravel2: ", err)
		return []schema.Packet{}
	}
	return results
}

func (this *PacketDatabase) GetPacketByPacketId(projectName string, packetId string) *schema.Packet {
	result := &schema.Packet{}
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
