package database

import (
	"fmt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type FileDatabase struct {
	Database
}

func (this *FileDatabase) Init() *FileDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&File{})
	return this
}

func (this *FileDatabase) GetFileByProjectAndId(project string, id string) *File {
	result := &File{}
	if err := this.db.FirstWithCtx(
		this.ctx,
		bson.M{"project": project, "fileId": id},
		result,
	); err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}
