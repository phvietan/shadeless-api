package database

import (
	"fmt"
	"shadeless-api/main/libs/database/schema"

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
	this.db = mgm.Coll(&schema.Project{})
	return this
}

func (this *ProjectDatabase) GetProjects() []schema.Project {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	results := []schema.Project{}
	if err := this.db.SimpleFind(&results, bson.M{}, findOptions); err != nil {
		fmt.Println(err)
		return []schema.Project{}
	}
	return results
}

func (this *ProjectDatabase) GetOneProjectById(id primitive.ObjectID) *schema.Project {
	project := &schema.Project{}
	if err := this.db.FindByID(id, project); err != nil {
		fmt.Println(err)
		return nil
	}
	return project
}

func (this *ProjectDatabase) GetOneProjectByName(name string) *schema.Project {
	project := &schema.Project{}
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

func (this *ProjectDatabase) UpdateProject(id primitive.ObjectID, project *schema.Project) error {
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
