package store

import "intro-rest/internal/app/interfaces"

type Store struct {
	entityRepo  interfaces.EntityRepo
	memDataRepo interfaces.MemDataRepo
}

func NewStore(entityRepo interfaces.EntityRepo, memDataRepo interfaces.MemDataRepo) interfaces.Store {
	return &Store{
		entityRepo:  entityRepo,
		memDataRepo: memDataRepo,
	}
}

func (s *Store) Entity() interfaces.EntityRepo {
	return s.entityRepo
}

func (s *Store) MemData() interfaces.MemDataRepo {
	return s.memDataRepo
}
