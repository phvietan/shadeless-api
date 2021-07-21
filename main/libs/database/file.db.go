package database

import (
	"github.com/kamva/mgm/v3"
)

type File struct {
	mgm.DefaultModel `bson:",inline"`

	Project string
	FileId  string
}

func NewFile(project string, fileId string) *File {
	return &File{
		Project: project,
		FileId:  fileId,
	}
}
