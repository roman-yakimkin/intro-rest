package memdatarepo

import (
	"intro-rest/internal/app/interfaces"
	"intro-rest/internal/app/models"
	"intro-rest/internal/app/services/configmanager"
	"os"
	"os/signal"
	"time"
)

type MemDataRepo struct {
	entities map[string]models.Entity
	er       interfaces.EntityRepo
	config   *configmanager.Config
}

func NewMemDataRepo(er interfaces.EntityRepo, config *configmanager.Config) *MemDataRepo {
	return &MemDataRepo{
		entities: make(map[string]models.Entity),
		er:       er,
		config:   config,
	}
}

func (r *MemDataRepo) Init() error {
	if err := r.copyAllFromDB(); err != nil {
		return err
	}
	ticker := time.NewTicker(time.Duration(r.config.CopyDataIntervalMS) * time.Millisecond)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		for {
			select {
			case <-signalCh:
				return
			case <-ticker.C:
				err := r.GetFromEntityRepo()
				if err != nil {
					return
				}
			}
		}
	}()
	return nil
}

func (r *MemDataRepo) GetFromEntityRepo() error {
	ids := r.er.GetDeleted()
	for _, id := range ids {
		delete(r.entities, id)
		r.er.CleanChanged(id)
	}
	insEntities, err := r.er.GetInsUpdated(interfaces.EntityInserted)
	if err != nil {
		return err
	}
	for _, entity := range insEntities {
		r.entities[entity.ID] = entity
		r.er.CleanChanged(entity.ID)
	}
	updEntities, err := r.er.GetInsUpdated(interfaces.EntityUpdated)
	if err != nil {
		return err
	}
	for _, entity := range updEntities {
		r.entities[entity.ID] = entity
		r.er.CleanChanged(entity.ID)
	}
	return nil
}

func (r *MemDataRepo) GetAll() ([]models.Entity, error) {
	result := make([]models.Entity, 0, len(r.entities))
	for _, entity := range r.entities {
		result = append(result, entity)
	}
	return result, nil
}

func (r *MemDataRepo) copyAllFromDB() error {
	entities, err := r.er.GetAll()
	if err != nil {
		return err
	}
	for _, entity := range entities {
		r.entities[entity.ID] = entity
	}
	return nil
}
