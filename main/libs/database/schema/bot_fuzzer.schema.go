package schema

import (
	"github.com/kamva/mgm/v3"
)

type BotFuzzer struct {
	mgm.DefaultModel `bson:",inline"`

	Project      string `json:"project" bson:"project"`
	Running      bool   `json:"running" bson:"running"`
	SleepRequest int    `json:"sleepRequest" bson:"sleepRequest"` // Sleep between requests in ms
	AsyncRequest int    `json:"asyncRequest" bson:"asyncRequest"` // Number of async requests
	Timeout      int    `json:"timeout" bson:"timeout"`           // ms
}

func NewBotFuzzer(project string) *BotFuzzer {
	return &BotFuzzer{
		Project:      project,
		Running:      false,
		SleepRequest: 0,     // Quite fast, does not sleep
		AsyncRequest: 5,     // Async 5 requests at a time
		Timeout:      10000, // ms
	}
}
