package models

import "github.com/google/uuid"

type (
	FoundingRange interface {
		ID() uuid.UUID
		Low() int
		High() int
	}

	CoFoundingRange interface {
		FoundingRange
	}

	foundingRange struct {
		id   uuid.UUID
		low  int
		high int
	}
)

func (f foundingRange) ID() uuid.UUID {
	return f.id
}

func (f foundingRange) Low() int {
	return f.low
}

func (f foundingRange) High() int {
	return f.high
}

func NewRange(id uuid.UUID, low int, high int) FoundingRange {
	return &foundingRange{
		id:   id,
		low:  low,
		high: high,
	}
}
