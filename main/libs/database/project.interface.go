package database

import (
	"shadeless-api/main/libs/database/schema"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProjectDatabase interface {
	IDatabase

	Init() *ProjectDatabase

	GetProjects() []schema.Project
	GetOneProjectById(id primitive.ObjectID) *schema.Project
	GetOneProjectByName(name string) *schema.Project

	UpdateProject(id primitive.ObjectID, project *schema.Project) error
	UpdateProjectStatus(id primitive.ObjectID, newStatus string) error
}
