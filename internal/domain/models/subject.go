package models

import "github.com/google/uuid"

type (
	Subject interface {
		ID() uuid.UUID
		Name() string
		Event() uuid.UUID
	}

	subject struct {
		id    uuid.UUID
		name  string
		event uuid.UUID
	}
)

func NewSubject(name string, event uuid.UUID) Subject {
	return &subject{
		id:    uuid.New(),
		name:  name,
		event: event,
	}
}

func (s *subject) ID() uuid.UUID {
	return s.id
}

func (s *subject) Name() string {
	return s.name
}

func (s *subject) Event() uuid.UUID {
	return s.event
}
