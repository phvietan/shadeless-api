package database

import (
	"fmt"
	"os"
	"shadeless-api/main/libs/database/schema"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IBotFuzzerDatabase interface {
	IDatabase

	Init() *BotFuzzerDatabase
}

type BotFuzzerDatabase struct {
	Database
}

func (this *BotFuzzerDatabase) Init() *BotFuzzerDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&schema.BotFuzzer{})
	mod := mongo.IndexModel{
		Keys:    bson.D{{"project", 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := this.db.Indexes().CreateOne(this.ctx, mod)
	if err != nil {
		fmt.Println("Error when creating index, ", err)
		os.Exit(0)
	}
	return this
}
