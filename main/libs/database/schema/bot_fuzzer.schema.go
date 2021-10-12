package schema

import (
	"github.com/kamva/mgm/v3"
)

type BotFuzzer struct {
	mgm.DefaultModel `bson:",inline"`

	Project      string `bson:"project"`
	Running      bool   `bson:"running"`
	SleepRequest int    `bson:"sleepRequest"` // Sleep between requests in ms
	AsyncRequest int    `bson:"asyncRequest"` // Number of async requests
}

func NewBotFuzzer(project string) *BotPath {
	return &BotPath{
		Project:      project,
		Running:      false,
		SleepRequest: 0, // Quite fast,
		AsyncRequest: 5, // Async 5 requests at a time
	}
}
