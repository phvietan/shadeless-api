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

type IUserDatabase interface {
	IDatabase
	Init() *UserDatabase
	GetUsers(project string) []schema.User
	GetUserByProjectAndCodename(project string, codename string) *schema.User
	Upsert(project string, codeName string)
}

type UserDatabase struct {
	Database
}

func (this *UserDatabase) Init() *UserDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&schema.User{})
	mod := mongo.IndexModel{
		Keys:    bson.D{{"project", 1}, {"codeName", 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := this.db.Indexes().CreateOne(this.ctx, mod)
	if err != nil {
		fmt.Println("Error when creating index, ", err)
		os.Exit(0)
	}
	return this
}

func (this *UserDatabase) GetUsers(project string) []schema.User {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	results := []schema.User{}
	if err := this.db.SimpleFind(&results, bson.M{
		"project": project,
	}, findOptions); err != nil {
		fmt.Println(err)
		return []schema.User{}
	}
	return results
}

func (this *UserDatabase) GetUserByProjectAndCodename(project string, codeName string) *schema.User {
	user := &schema.User{}
	if err := this.db.FirstWithCtx(
		this.ctx,
		bson.M{"project": project, "codeName": codeName},
		user,
	); err != nil {
		return nil
	}
	return user
}

func (this *UserDatabase) Upsert(project string, codeName string) {
	user := this.GetUserByProjectAndCodename(project, codeName)
	if user == nil {
		this.Insert(schema.NewUser(project, codeName))
	}
}
