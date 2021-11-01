package database

import (
	"errors"
	"fmt"
	"os"
	"shadeless-api/main/libs/database/schema"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IBotPathDatabase interface {
	IDatabase

	Init() *BotPathDatabase
	GetBotPathByProject(projectName string) *schema.BotPath
	PutBotPathByProject(id primitive.ObjectID, newBotPath *schema.BotPath) error
	SwitchRun(botPath *schema.BotPath) error
}

type BotPathDatabase struct {
	Database
}

func (this *BotPathDatabase) Init() *BotPathDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&schema.BotPath{})
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

func (this *BotPathDatabase) GetBotPathByProject(projectName string) *schema.BotPath {
	result := &schema.BotPath{}
	if err := this.db.FirstWithCtx(
		this.ctx,
		bson.M{"project": projectName},
		result,
	); err != nil {
		fmt.Errorf("%v", err)
		return nil
	}
	return result
}

func (this *BotPathDatabase) PutBotPathByProject(id primitive.ObjectID, newBotPath *schema.BotPath) error {
	if newBotPath == nil {
		return errors.New("Error: BotPath is undefined")
	}
	updated := bson.M{
		"sleepRequest": newBotPath.SleepRequest,
		"asyncRequest": newBotPath.AsyncRequest,
		"timeout":      newBotPath.Timeout,
	}

	if _, err := this.db.UpdateByID(this.ctx, id, bson.D{{"$set", updated}}); err != nil {
		return err
	}
	return nil
}

func (this *BotPathDatabase) SwitchRun(botPath *schema.BotPath) error {
	if botPath == nil {
		return errors.New("Error: BotPath is undefined")
	}
	updated := bson.M{
		"running": !botPath.Running,
	}

	if _, err := this.db.UpdateByID(this.ctx, botPath.ID, bson.D{{"$set", updated}}); err != nil {
		return err
	}
	return nil
}
