package interfaces

import "intro-rest/internal/app/models"

type MemDataRepo interface {
	Init() error
	GetFromEntityRepo() error
	GetAll() ([]models.Entity, error)
}
