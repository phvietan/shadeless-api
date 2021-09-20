package database

import (
	"shadeless-api/main/libs"
	"shadeless-api/main/libs/database/schema"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestFileCreateNilFile(t *testing.T) {
	var dbInstance IFileDatabase = new(FileDatabase).Init()
	err := dbInstance.Insert(nil)
	assert.Equal(t, err.Error(), "Object to insert to mongo is nil")
}

func TestFileCreateAndQueryFile(t *testing.T) {
	var dbInstance IFileDatabase = new(FileDatabase).Init()
	defer dbInstance.ClearCollection()
	for i := 0; i < 10; i++ {
		projectName := libs.RandomString(32)
		fileId := libs.RandomString(32)
		f := schema.NewFile(projectName, fileId)
		err := dbInstance.Insert(f)
		assert.Equal(t, err, nil)
		fDB := dbInstance.GetFileByProjectAndId(projectName, fileId)
		assert.NotEqual(t, fDB, nil)
	}
}

func TestQueryUnknownFile(t *testing.T) {
	var dbInstance IFileDatabase = new(FileDatabase).Init()
	defer dbInstance.ClearCollection()
	for i := 0; i < 10; i++ {
		projectName := libs.RandomString(32)
		fileId := libs.RandomString(32)
		fDB := dbInstance.GetFileByProjectAndId(projectName, fileId)
		assert.Equal(t, fDB, nil)
	}
}

func TestCreateThenDeleteFiles(t *testing.T) {
	var dbInstance IFileDatabase = new(FileDatabase).Init()
	defer dbInstance.ClearCollection()
	arrName := []string{}
	arrId := []string{}
	for i := 0; i < 10; i++ {
		projectName := libs.RandomString(64)
		fileId := libs.RandomString(64)
		err := dbInstance.Insert(schema.NewFile(projectName, fileId))
		assert.Equal(t, err, nil)
		arrName = append(arrName, projectName)
		arrId = append(arrId, fileId)
	}

	for i := 0; i < 10; i++ {
		projectName := arrName[i]
		fileId := arrId[i]
		file := dbInstance.GetFileByProjectAndId(projectName, fileId)
		assert.NotEqual(t, file, nil)
		err := dbInstance.DeleteByOneProperty("project", projectName)
		assert.Equal(t, err, nil)
		file2 := dbInstance.GetFileByProjectAndId(projectName, fileId)
		assert.Equal(t, file2, nil)
	}
}

func TestDeleteFileNonExistingProject(t *testing.T) {
	var dbInstance IFileDatabase = new(FileDatabase).Init()
	err := dbInstance.DeleteByOneProperty("project", "aaa")
	assert.Equal(t, err, nil)
}
