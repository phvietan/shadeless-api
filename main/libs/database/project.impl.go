package database

import (
	"errors"
	"fmt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectDatabase struct {
	Database
}

func (this *ProjectDatabase) Init() *ProjectDatabase {
	this.ctx = mgm.Ctx()
	this.db = mgm.Coll(&Project{})
	return this
}

func (this *ProjectDatabase) CreateProject(project *Project) error {
	if project == nil {
		return errors.New("Project object is nil")
	}
	if err := this.db.Create(project); err != nil {
		fmt.Println(err)
	}
	return nil
}

func (this *ProjectDatabase) GetProjects() []Project {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	results := []Project{}
	if err := this.db.SimpleFind(&results, bson.M{}, findOptions); err != nil {
		fmt.Println(err.Error())
		return []Project{}
	}
	return results
}

func (this *ProjectDatabase) GetOneProjectById(id primitive.ObjectID) *Project {
	project := &Project{}
	if err := this.db.FindByID(id, project); err != nil {
		fmt.Errorf("%v", err)
		return nil
	}
	return project
}

func (this *ProjectDatabase) GetOneProjectByName(name string) *Project {
	project := &Project{}
	if err := this.db.First(bson.M{"name": name}, project); err != nil {
		fmt.Errorf("%v", err)
		return nil
	}
	return project
}

func (this *ProjectDatabase) GetNumberDocumentsByProject(name string) int64 {
	num, err := this.db.CountDocuments(this.ctx, bson.M{"project": name})
	if err != nil {
		fmt.Errorf("%v", err)
		return -1
	}
	return num
}

func (this *ProjectDatabase) UpdateProject(id primitive.ObjectID, project *Project) error {
	updated := bson.M{}
	if project.Name != "" {
		updated["name"] = project.Name
	}
	if project.Description != "" {
		updated["description"] = project.Description
	}
	if project.Status != "" {
		updated["status"] = project.Status
	}

	if _, err := this.db.UpdateByID(this.ctx, id, bson.D{{"$set", updated}}); err != nil {
		return err
	}
	return nil
}

func (this *ProjectDatabase) DeleteProject(id primitive.ObjectID) error {
	if _, err := this.db.DeleteOne(this.ctx, bson.M{"_id": id}); err != nil {
		return err
	}
	return nil
}
