package models

import "github.com/google/uuid"

type (
	Subject interface {
		ID() uuid.UUID
		Name() string
	}

	subject struct {
		id   uuid.UUID
		name string
	}
)

func NewSubject(name string) Subject {
	return &subject{
		id:   uuid.New(),
		name: name,
	}
}

func (s *subject) ID() uuid.UUID {
	return s.id
}

func (s *subject) Name() string {
	return s.name
}
