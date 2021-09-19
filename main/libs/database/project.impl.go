package database

import (
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

func (this *ProjectDatabase) GetProjects() []Project {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	results := []Project{}
	if err := this.db.SimpleFind(&results, bson.M{}, findOptions); err != nil {
		fmt.Println(err)
		return []Project{}
	}
	return results
}

func (this *ProjectDatabase) GetOneProjectById(id primitive.ObjectID) *Project {
	project := &Project{}
	if err := this.db.FindByID(id, project); err != nil {
		fmt.Println(err)
		return nil
	}
	return project
}

func (this *ProjectDatabase) GetOneProjectByName(name string) *Project {
	project := &Project{}
	if err := this.db.FirstWithCtx(
		this.ctx,
		bson.M{"name": name},
		project,
	); err != nil {
		fmt.Println(err)
		return nil
	}
	return project
}

func (this *ProjectDatabase) UpdateProject(id primitive.ObjectID, project *Project) error {
	updated := bson.M{
		"name":        project.Name,
		"description": project.Description,
		"blacklist":   project.Blacklist,
	}

	if _, err := this.db.UpdateByID(this.ctx, id, bson.D{{"$set", updated}}); err != nil {
		return err
	}
	return nil
}

func (this *ProjectDatabase) UpdateProjectStatus(id primitive.ObjectID, newStatus string) error {
	if _, err := this.db.UpdateByID(this.ctx, id,
		bson.D{{"$set", bson.M{
			"status": newStatus,
		}}},
	); err != nil {
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
