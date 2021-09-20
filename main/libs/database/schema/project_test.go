package schema

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNewProject(t *testing.T) {
	p := NewProject()
	assert.Equal(t, len(p.Blacklist), 0)
	assert.Equal(t, p.Status, ProjectStatusTodo)
}
