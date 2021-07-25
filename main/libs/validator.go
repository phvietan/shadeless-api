package libs

import (
	"errors"
	"regexp"
)

type validator struct{}

func (this *validator) ValidateProjectName(name string) error {
	if match, err := regexp.MatchString("^[a-zA-Z0-9]+$", name); err != nil || !match {
		if err != nil {
			return err
		}
		return errors.New("Project name should match ^[a-zA-Z0-9]+$")
	}
	return nil
}

var Validator = new(validator)
