package models

import "github.com/google/uuid"

type (
	OrganizerLevel interface {
		ID() uuid.UUID
		Name() string
		Code() string
	}

	organizerLevel struct {
		id   uuid.UUID
		name string
		code string
	}
)

func NewOrganizerLevel(id uuid.UUID, name string, code string) OrganizerLevel {
	return &organizerLevel{
		id:   id,
		name: name,
		code: code,
	}
}

func (l *organizerLevel) ID() uuid.UUID {
	return l.id
}

func (l *organizerLevel) Name() string {
	return l.name
}

func (l *organizerLevel) Code() string {
	return l.code
}
