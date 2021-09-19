package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type IProjectDatabase interface {
	IDatabase

	Init() *ProjectDatabase

	GetProjects() []Project
	GetOneProjectById(id primitive.ObjectID) *Project
	GetOneProjectByName(name string) *Project

	UpdateProject(id primitive.ObjectID, project *Project) error
	UpdateProjectStatus(id primitive.ObjectID, newStatus string) error
	DeleteProject(id primitive.ObjectID) error

	ClearCollection()
}
