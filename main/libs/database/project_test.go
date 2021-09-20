package database

import (
	"shadeless-api/main/libs"
	"shadeless-api/main/libs/database/schema"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestFilterBlacklistMatch(t *testing.T) {
	for i := 0; i < 20; i++ {
		blacklist := []schema.Blacklist{}
		blacklistValue, err := schema.NewBlacklist(libs.RandomString(32), schema.BlacklistValue)
		assert.Equal(t, err, nil)
		blacklist = append(blacklist, *blacklistValue)

		blacklistRegex, err := schema.NewBlacklist("^a", schema.BlacklistRegex)
		assert.Equal(t, err, nil)
		blacklist = append(blacklist, *blacklistRegex)

		origins := []string{
			blacklistValue.Value,
			libs.RandomString(32),
		}
		filtered := schema.FilterBlacklistMatch(origins, blacklist)
		assert.Equal(t, len(filtered) <= 1, true)
		if origins[1][0] == 'a' {
			assert.Equal(t, len(filtered) == 0, true)
		} else {
			assert.Equal(t, len(filtered) == 1, true)
		}
	}
}

func TestCreateProject(t *testing.T) {
	var dbInstance IProjectDatabase = new(ProjectDatabase).Init()
	defer dbInstance.ClearCollection()
	err := dbInstance.Insert(nil)
	assert.Equal(t, err.Error(), "Object to insert to mongo is nil")

	for i := 0; i < 10; i++ {
		newProject := schema.NewProject()
		newProject.Name = libs.RandomString(32)
		err := dbInstance.Insert(newProject)
		assert.Equal(t, err, nil)

		allProjects := dbInstance.GetProjects()
		assert.Equal(t, len(allProjects), i+1)
	}
}

func TestCreateAndGetAllProject(t *testing.T) {
	var dbInstance IProjectDatabase = new(ProjectDatabase).Init()
	defer dbInstance.ClearCollection()
	for i := 0; i < 10; i++ {
		newProject := schema.NewProject()
		newProject.Name = libs.RandomString(32)
		err := dbInstance.Insert(newProject)
		assert.Equal(t, err, nil)

		allProjects := dbInstance.GetProjects()
		assert.Equal(t, len(allProjects), i+1)
	}
}

func TestProjectQuery(t *testing.T) {
	var dbInstance IProjectDatabase = new(ProjectDatabase).Init()
	defer dbInstance.ClearCollection()
	for i := 0; i < 10; i++ {
		newProject := schema.NewProject()
		newProject.Name = libs.RandomString(32)
		err := dbInstance.Insert(newProject)
		assert.Equal(t, err, nil)

		dbProjectByName := dbInstance.GetOneProjectByName(newProject.Name)
		assert.NotEqual(t, dbProjectByName, nil)
		assert.Equal(t, dbProjectByName.Name, newProject.Name)

		dbProjectById := dbInstance.GetOneProjectById(dbProjectByName.ID)
		assert.NotEqual(t, dbProjectById, nil)
		assert.Equal(t, dbProjectById.Name, newProject.Name)

		dbNonExistProjectByName := dbInstance.GetOneProjectByName(newProject.Name + "dcm")
		assert.Equal(t, dbNonExistProjectByName, nil)

		dbNonExistProjectById := dbInstance.GetOneProjectById([12]byte{})
		assert.Equal(t, dbNonExistProjectById, nil)
	}
}

func TestProjectUpdate(t *testing.T) {
	var dbInstance IProjectDatabase = new(ProjectDatabase).Init()
	defer dbInstance.ClearCollection()
	for i := 0; i < 10; i++ {
		newProject := schema.NewProject()
		newProject.Name = libs.RandomString(32)
		newProject.Description = libs.RandomString(32)
		err := dbInstance.Insert(newProject)
		assert.Equal(t, err, nil)

		dbProjectByName := dbInstance.GetOneProjectByName(newProject.Name)
		assert.NotEqual(t, dbProjectByName, nil)
		assert.Equal(t, dbProjectByName.Name, newProject.Name)

		newName := libs.RandomString(32)
		newDescription := libs.RandomString(32)
		blacklist, err := schema.NewBlacklist("a", schema.BlacklistRegex)
		assert.Equal(t, err, nil)

		updateProject := schema.NewProject()
		updateProject.Name = newName
		updateProject.Description = newDescription
		updateProject.Status = schema.ProjectStatusDone
		updateProject.Blacklist = []schema.Blacklist{*blacklist}

		err = dbInstance.UpdateProject(dbProjectByName.ID, updateProject)
		assert.Equal(t, err, nil)

		dbProjectById := dbInstance.GetOneProjectById(dbProjectByName.ID)
		assert.NotEqual(t, dbProjectById, nil)

		assert.Equal(t, dbProjectById.Name, newName)
		assert.Equal(t, dbProjectById.Description, newDescription)
		assert.Equal(t, len(dbProjectById.Blacklist[0].Value), 1)
		assert.Equal(t, dbProjectById.Blacklist[0].Value, "a")
		assert.Equal(t, dbProjectById.Blacklist[0].Type, schema.BlacklistRegex)
	}
}

func TestProjectDelete(t *testing.T) {
	var dbInstance IProjectDatabase = new(ProjectDatabase).Init()
	defer dbInstance.ClearCollection()
	err := dbInstance.DeleteProject([12]byte{})
	assert.Equal(t, err, nil)

	for i := 0; i < 10; i++ {
		newProject := schema.NewProject()
		newProject.Name = libs.RandomString(32)
		err = dbInstance.Insert(newProject)
		assert.Equal(t, err, nil)

		allProjects := dbInstance.GetProjects()
		assert.Equal(t, len(allProjects), 1)

		dbProjectByName := dbInstance.GetOneProjectByName(newProject.Name)
		assert.NotEqual(t, dbProjectByName, nil)
		assert.Equal(t, dbProjectByName.Name, newProject.Name)

		err = dbInstance.DeleteProject(dbProjectByName.ID)
		assert.Equal(t, err, nil)

		allProjects = dbInstance.GetProjects()
		assert.Equal(t, len(allProjects), 0)
	}
}
