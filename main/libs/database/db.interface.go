package database

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDatabase interface {
	Insert(row mgm.Model) error
	ClearCollection()
	UpdateOneProperty(propertyKey string, propertyOldValue interface{}, propertyNewValue interface{}) error
	DeleteByOneProperty(propertyKey string, propertyValue interface{}) error
	DeleteById(id primitive.ObjectID) error
}
