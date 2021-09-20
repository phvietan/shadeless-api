package schema

import "github.com/kamva/mgm/v3"

type Note struct {
	mgm.DefaultModel `bson:",inline"`

	UserId          string  `json:"userId" bson:"userId"`
	Tags            string  `json:"tags"`
	Description     string  `json:"description"`
	RequestPacketId string  `json:"requestPacketId" bson:"requestPacketId"`
	Replies         []Reply `json:"replies" bson:"replies"`
}

type Reply struct {
	NoteId      string `json:"noteId"`
	UserId      string `json:"userId" bson:"userId"`
	Description string `json:"description"`
}
