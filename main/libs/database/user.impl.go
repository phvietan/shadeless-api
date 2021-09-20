package database

import (
	"fmt"
	"shadeless-api/main/libs/database/schema"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDatabase struct {
	Database
}

func (this *UserDatabase) Init() *UserDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&schema.User{})
	return this
}

func (this *UserDatabase) GetUsers() []schema.User {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	results := []schema.User{}
	if err := this.db.SimpleFind(&results, bson.M{}, findOptions); err != nil {
		fmt.Println(err)
		return []schema.User{}
	}
	return results
}
