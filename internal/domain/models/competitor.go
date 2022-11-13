package models

import "github.com/google/uuid"

type (
	Competitor interface {
		ID() uuid.UUID
		Name() string
	}

	competitor struct {
		id   uuid.UUID
		name string
	}
)

func NewCompetitor(name string) Competitor {
	return &competitor{
		id:   uuid.New(),
		name: name,
	}
}

func (c *competitor) ID() uuid.UUID {
	return c.id
}

func (c *competitor) Name() string {
	return c.name
}
