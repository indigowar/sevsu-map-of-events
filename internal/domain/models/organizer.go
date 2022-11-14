package models

import (
	"github.com/google/uuid"
)

type (
	Organizer interface {
		ID() uuid.UUID
		Name() string
		Logo() string
		Level() uuid.UUID
	}

	organizer struct {
		id    uuid.UUID
		name  string
		logo  string
		level uuid.UUID
	}
)

func (o organizer) ID() uuid.UUID {
	return o.id
}

func (o organizer) Name() string {
	return o.name
}

func (o organizer) Logo() string {
	return o.logo
}

func (o organizer) Level() uuid.UUID {
	return o.level
}

func NewOrganizer(id uuid.UUID, name, logo string, level uuid.UUID) Organizer {
	return &organizer{
		id:    id,
		name:  name,
		logo:  logo,
		level: level,
	}
}
