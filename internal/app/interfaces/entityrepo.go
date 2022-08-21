package interfaces

import "intro-rest/internal/app/models"

const (
	EntityDeleted = iota
	EntityInserted
	EntityUpdated
)

type EntityRepo interface {
	Init() error
	Save(*models.Entity) (string, error)
	Delete(string) error
	GetByID(string) (*models.Entity, error)
	GetAll() ([]models.Entity, error)
	GetInsUpdated(int) ([]models.Entity, error)
	GetDeleted() []string
	CleanChanged(string)
}
