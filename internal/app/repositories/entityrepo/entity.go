package entityrepo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"intro-rest/internal/app/models"
	"time"
)

type Entity struct {
	ID      primitive.ObjectID `bson:"_id, omitempty"`
	Name    string             `bson:"name"`
	Created time.Time          `bson:"created"`
}

func (e *Entity) Export() *models.Entity {
	var entity models.Entity
	entity.ID = e.ID.Hex()
	entity.Name = e.Name
	entity.Created = e.Created
	return &entity
}

func (e *Entity) Import(entity *models.Entity) error {
	var err error
	e.ID = primitive.ObjectID{}
	if entity.ID != "" {
		e.ID, err = primitive.ObjectIDFromHex(entity.ID)
		if err != nil {
			return err
		}
	}
	e.Name = entity.Name
	e.Created = entity.Created
	return nil
}
