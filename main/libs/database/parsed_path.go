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

type IParsedPathDatabase interface {
	IDatabase

	Init() *ParsedPathDatabase
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
