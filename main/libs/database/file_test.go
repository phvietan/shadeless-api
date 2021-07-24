package database

import (
	"shadeless-api/main/libs"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestFileDb(t *testing.T) {
	projectName := libs.RandomString(32)
	fileId := libs.RandomString(32)
	f := NewFile(projectName, fileId)
	assert.Equal(t, projectName, f.Project)
	assert.Equal(t, fileId, f.FileId)
}

func TestFileCreateNilFile(t *testing.T) {
	var dbInstance IFileDatabase = new(FileDatabase).Init()
	err := dbInstance.CreateFile(nil)
	assert.Equal(t, err.Error(), "Project object is nil")
}

func TestFileCreateAndQueryFile(t *testing.T) {
	var dbInstance IFileDatabase = new(FileDatabase).Init()
	defer dbInstance.ClearCollection()
	for i := 0; i < 10; i++ {
		projectName := libs.RandomString(32)
		fileId := libs.RandomString(32)
		f := NewFile(projectName, fileId)
		err := dbInstance.CreateFile(f)
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
		err := dbInstance.CreateFile(NewFile(projectName, fileId))
		assert.Equal(t, err, nil)
		arrName = append(arrName, projectName)
		arrId = append(arrId, fileId)
	}

	for i := 0; i < 10; i++ {
		projectName := arrName[i]
		fileId := arrId[i]
		file := dbInstance.GetFileByProjectAndId(projectName, fileId)
		assert.NotEqual(t, file, nil)
		err := dbInstance.DeleteFilesByProjectName(projectName)
		assert.Equal(t, err, nil)
		file2 := dbInstance.GetFileByProjectAndId(projectName, fileId)
		assert.Equal(t, file2, nil)
	}
}

func TestDeleteFileNonExistingProject(t *testing.T) {
	var dbInstance IFileDatabase = new(FileDatabase).Init()
	err := dbInstance.DeleteFilesByProjectName("aaa")
	assert.Equal(t, err, nil)
}
