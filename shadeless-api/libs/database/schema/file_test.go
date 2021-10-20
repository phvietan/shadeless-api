package schema

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
