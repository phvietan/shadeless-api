package database

import "shadeless-api/main/libs/database/schema"

type IFileDatabase interface {
	IDatabase
	Init() *FileDatabase
	GetFileByProjectAndId(project string, id string) *schema.File
}
