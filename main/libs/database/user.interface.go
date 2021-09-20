package database

import (
	"shadeless-api/main/libs/database/schema"
)

type IUserDatabase interface {
	IDatabase
	Init() *UserDatabase
	GetUsers() []schema.User
}
