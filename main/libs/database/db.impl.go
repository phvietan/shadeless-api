package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"shadeless-api/main/config"

	"github.com/benweissmann/memongo"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	ctx context.Context
	db  *mgm.Collection
}

func (this *Database) Insert(row mgm.Model) error {
	if row == nil {
		return errors.New("Object to insert to mongo is nil")
	}
	return this.db.CreateWithCtx(this.ctx, row)
}

func (this *Database) UpdateOneProperty(propertyKey string, propertyOldValue interface{}, propertyNewValue interface{}) error {
	_, err := this.db.UpdateMany(
		this.ctx,
		bson.M{propertyKey: propertyOldValue},
		bson.D{{"$set", bson.M{propertyKey: propertyNewValue}}},
	)
	return err
}

func (this *Database) DeleteByOneProperty(propertyKey string, propertyValue interface{}) error {
	if _, err := this.db.DeleteMany(this.ctx, bson.M{propertyKey: propertyValue}); err != nil {
		return err
	}
	return nil
}

func (this *Database) ClearCollection() {
	if config.GetInstance().GetEnvironment() != "test" {
		fmt.Println("Only allow to run this function in test mode")
		return
	}
	this.db.Drop(this.ctx)
}

var mongoServer *memongo.Server = nil

func init() {
	if config.GetInstance().GetEnvironment() == "test" {
		fmt.Println("Test suite: creating mongo memory database")
		var err error
		mongoServer, err = memongo.Start("4.0.5")
		if err != nil {
			panic(err)
		}
		databaseUrl := mongoServer.URIWithRandomDB()
		os.Setenv("DATABASE_URL", databaseUrl)
		config.GetInstance().SetDatabaseUrl(databaseUrl)
	}
	databaseUrl := config.GetInstance().GetDatabaseUrl()
	err := mgm.SetDefaultConfig(
		nil,
		"shadeless",
		options.Client().ApplyURI(databaseUrl),
	)
	if err != nil {
		panic(err)
	}
}
