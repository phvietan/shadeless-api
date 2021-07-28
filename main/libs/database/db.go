package database

import (
	"context"
	"fmt"
	"os"
	"shadeless-api/main/config"

	"github.com/benweissmann/memongo"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	ctx context.Context
	db  *mgm.Collection
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
