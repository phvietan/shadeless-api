package database

import (
	"errors"

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
	return this.db.Create(file)
}

func (this *FileDatabase) GetFileByProjectAndId(project string, id string) *File {
	result := &File{}
	if err := this.db.First(bson.M{"project": project, "fileId": id}, result); err != nil {
		return nil
	}
	return result
}

func (this *FileDatabase) DeleteFilesByProjectName(projectName string) error {
	if _, err := this.db.DeleteMany(this.ctx, bson.M{"project": projectName}); err != nil {
		return err
	}
	return nil
}

func (this *FileDatabase) UpdateProjectName(oldName string, newName string) error {
	_, err := this.db.UpdateMany(
		this.ctx,
		bson.M{"project": oldName},
		bson.D{{"$set", bson.M{"project": newName}}},
	)
	return err
}
