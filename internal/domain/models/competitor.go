package models

import "github.com/google/uuid"

type Competitor struct {
	ID   uuid.UUID
	Name string
}
