package database

import (
	"github.com/kamva/mgm/v3"
)

const (
	ProjectStatusTodo    = "todo"
	ProjectStatusHacking = "hacking"
	ProjectStatusDone    = "done"
)

type Project struct {
	mgm.DefaultModel `bson:",inline"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func NewProject() *Project {
	return &Project{
		Status: ProjectStatusTodo,
	}
}
