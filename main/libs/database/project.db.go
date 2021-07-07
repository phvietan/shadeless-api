package database

import (
	"errors"
	"fmt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Project struct {
	mgm.DefaultModel `bson:",inline"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

const (
	ProjectStatusTodo    = "todo"
	ProjectStatusHacking = "hacking"
	ProjectStatusDone    = "done"
)

func NewProject() *Project {
	return &Project{
		Status: ProjectStatusTodo,
	}
}

func CreateProject(project *Project) error {
	if project == nil {
		return errors.New("Project object is nil")
	}
	err := mgm.Coll(project).Create(project)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func GetProjects() []Project {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	results := []Project{}
	coll := mgm.Coll(&Project{})
	err := coll.SimpleFind(&results, bson.M{}, findOptions)
	if err != nil {
		fmt.Println(err.Error())
		return []Project{}
	}
	return results
}

func GetOneProjectById(id primitive.ObjectID) *Project {
	project := &Project{}
	err := mgm.Coll(project).FindByID(id, project)
	if err != nil {
		fmt.Errorf("%v", err)
		return nil
	}
	return project
}

func GetOneProjectByName(name string) *Project {
	project := &Project{}
	err := mgm.Coll(project).First(bson.M{"name": name}, project)
	if err != nil {
		fmt.Errorf("%v", err)
		return nil
	}
	return project
}

func UpdateProject(id primitive.ObjectID, project *Project) error {
	ctx := mgm.Ctx()

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

	coll := mgm.Coll(&Project{})
	_, err := coll.UpdateByID(ctx, id, bson.D{
		{"$set", updated},
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteProject(id primitive.ObjectID) error {
	ctx := mgm.Ctx()

	coll := mgm.Coll(&Project{})
	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
