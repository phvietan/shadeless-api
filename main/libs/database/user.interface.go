package database

import (
	"shadeless-api/main/libs/database/schema"
)

type IUserDatabase interface {
	IDatabase
	Init() *UserDatabase
	GetUsers(project string) []schema.User
	GetUserByProjectAndCodename(project string, codename string) *schema.User
	Upsert(project string, codeName string)
}
