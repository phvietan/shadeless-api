package database

import (
	"context"
	"shadeless-api/main/config"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	ctx context.Context
	db  *mgm.Collection
}

func init() {
	err := mgm.SetDefaultConfig(
		nil,
		"shadeless",
		options.Client().ApplyURI(config.GetInstance().GetDatabaseUrl()),
	)
	if err != nil {
		panic(err)
	}
}
