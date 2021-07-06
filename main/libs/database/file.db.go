package database

import (
	"errors"
	"fmt"

	"github.com/kamva/mgm/v3"
)

type File struct {
	mgm.DefaultModel `bson:",inline"`

	Project string
	FileId  string
}

func CreateFile(file *File) error {
	if file == nil {
		return errors.New("Project object is nil")
	}
	err := mgm.Coll(file).Create(file)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
