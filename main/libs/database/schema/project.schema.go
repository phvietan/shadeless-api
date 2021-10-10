package schema

import (
	"errors"
	"regexp"

	"github.com/kamva/mgm/v3"
)

type Project struct {
	mgm.DefaultModel `bson:",inline"`

	Name        string      `json:"name"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	Blacklist   []Blacklist `json:"blacklist"`
}

func NewProject() *Project {
	return &Project{
		Status:    ProjectStatusTodo,
		Blacklist: make([]Blacklist, 0),
	}
}

func (this *Project) validateProjectBlacklist() error {
	cnt := 0
	for _, bl := range this.Blacklist {
		if bl.Type == BlacklistRegex {
			cnt += 1
		}
	}
	if cnt > 1 {
		return errors.New("Blacklist should have 1 regex only")
	}
	return nil
}
func (this *Project) validateProjectName() error {
	if match, err := regexp.MatchString("^[a-zA-Z0-9]+$", this.Name); err != nil || !match {
		if err != nil {
			return err
		}
		return errors.New("Project name should match ^[a-zA-Z0-9]+$")
	}
	return nil
}

func (this *Project) Validate() error {
	if err := this.validateProjectName(); err != nil {
		return err
	}
	return this.validateProjectBlacklist()
}
