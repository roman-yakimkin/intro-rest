package models

import "time"

type Entity struct {
	ID      string
	Name    string
	Created time.Time
}

func (e *Entity) BeforeSave() {
	if e.Created.IsZero() {
		e.Created = time.Now()
	}
}
