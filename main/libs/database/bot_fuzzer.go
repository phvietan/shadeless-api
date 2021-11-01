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

type IBotFuzzerDatabase interface {
	IDatabase

	Init() *BotFuzzerDatabase
	GetBotFuzzerByProject(projectName string) *schema.BotFuzzer
	PutBotFuzzerByProject(id primitive.ObjectID, newBotFuzzer *schema.BotFuzzer) error
	SwitchRun(BotFuzzer *schema.BotFuzzer) error
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

func (this *BotFuzzerDatabase) GetBotFuzzerByProject(projectName string) *schema.BotFuzzer {
	result := &schema.BotFuzzer{}
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

func (this *BotFuzzerDatabase) PutBotFuzzerByProject(id primitive.ObjectID, newBotFuzzer *schema.BotFuzzer) error {
	if newBotFuzzer == nil {
		return errors.New("Error: BotFuzzer is undefined")
	}
	updated := bson.M{
		"sleepRequest": newBotFuzzer.SleepRequest,
		"asyncRequest": newBotFuzzer.AsyncRequest,
		"timeout":      newBotFuzzer.Timeout,
	}

	if _, err := this.db.UpdateByID(this.ctx, id, bson.D{{"$set", updated}}); err != nil {
		return err
	}
	return nil
}

func (this *BotFuzzerDatabase) SwitchRun(BotFuzzer *schema.BotFuzzer) error {
	if BotFuzzer == nil {
		return errors.New("Error: BotFuzzer is undefined")
	}
	updated := bson.M{
		"running": !BotFuzzer.Running,
	}

	if _, err := this.db.UpdateByID(this.ctx, BotFuzzer.ID, bson.D{{"$set", updated}}); err != nil {
		return err
	}
	return nil
}
