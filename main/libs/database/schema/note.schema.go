package schema

import "github.com/kamva/mgm/v3"

type Note struct {
	mgm.DefaultModel `bson:",inline"`

	Project         string `json:"project" bson:"project"`
	CodeName        string `json:"codeName" bson:"codeName"`
	RequestPacketId string `json:"requestPacketId" bson:"requestPacketId"`
	Tags            string `json:"tags"`
	Description     string `json:"description"`
}

type Reply struct {
	CodeName    string `json:"codeName" bson:"codeName"`
	Description string `json:"description"`
}

func NewNote() *Note {
	return &Note{
		CodeName:        "",
		RequestPacketId: "",
		Tags:            "",
		Description:     "",
	}
}
