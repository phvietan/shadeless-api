package database

import (
	"errors"
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

func (this *FileDatabase) CreateFile(file *File) error {
	if file == nil {
		return errors.New("Project object is nil")
	}
	if err := this.db.Create(file); err != nil {
		fmt.Errorf("%v", err)
		return err
	}
	return nil
}

func (this *FileDatabase) DeleteFilesByProjectName(projectName string) error {
	if _, err := this.db.DeleteMany(this.ctx, bson.M{"project": projectName}); err != nil {
		return err
	}
	return nil
}
