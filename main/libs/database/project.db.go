package database

import (
	"errors"
	"fmt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
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
	ctx := mgm.Ctx()
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	coll := mgm.Coll(&Project{})
	cursor, err := coll.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		fmt.Println(err.Error())
		return []Project{}
	}

	results := []Project{}
	cursor.All(ctx, &results)
	return results
}
