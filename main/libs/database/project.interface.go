package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type IProjectDatabase interface {
	Init() *ProjectDatabase
	CreateProject(project *Project) error

	GetProjects() []Project
	GetOneProjectById(id primitive.ObjectID) *Project
	GetOneProjectByName(name string) *Project
	GetNumberDocumentsByProject(name string) int64

	UpdateProject(id primitive.ObjectID, project *Project) error
	DeleteProject(id primitive.ObjectID) error
}
