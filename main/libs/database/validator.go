package database

import (
	"errors"
	"regexp"
)

type validator struct{}

func (this *validator) ValidateProjectBlacklist(blacklist []Blacklist) error {
	cnt := 0
	for _, bl := range blacklist {
		if bl.Type == BlacklistRegex {
			cnt += 1
		}
	}
	if cnt > 1 {
		return errors.New("Blacklist should have 1 regex only")
	}
	return nil
}
func (this *validator) ValidateProjectName(name string) error {
	if match, err := regexp.MatchString("^[a-zA-Z0-9]*$", name); err != nil || !match {
		if err != nil {
			return err
		}
		return errors.New("Project name should match ^[a-zA-Z0-9]+$")
	}
	return nil
}

var Validator = new(validator)
