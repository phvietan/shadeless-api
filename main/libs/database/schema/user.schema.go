package schema

import (
	"shadeless-api/main/libs"

	"github.com/kamva/mgm/v3"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`

	CodeName string `json:"codeName" bson:"codeName"`
	Project  string `json:"project"`
	Color    string `json:"color"`
}

func NewUser(project string, codeName string) *User {
	user := new(User)
	user.CodeName = codeName
	user.Project = project
	user.Color = "#" + libs.RandomString(6)
	return user
}
