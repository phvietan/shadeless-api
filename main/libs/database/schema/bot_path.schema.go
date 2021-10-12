package schema

import (
	"github.com/kamva/mgm/v3"
)

type BotPath struct {
	mgm.DefaultModel `bson:",inline"`

	Project      string `json:"project" bson:"project"`
	Running      bool   `json:"running" bson:"running"`
	SleepRequest int    `json:"sleepRequest" bson:"sleepRequest"` // Sleep between requests in ms
	AsyncRequest int    `json:"asyncRequest" bson:"asyncRequest"` // Number of async requests
}

func NewBotPath(project string) *BotPath {
	return &BotPath{
		Project:      project,
		Running:      false,
		SleepRequest: 0, // Quite fast,
		AsyncRequest: 5, // Async 5 requests at a time
	}
}