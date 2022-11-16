package models

import (
	"github.com/google/uuid"
)

type Organizer struct {
	ID    uuid.UUID
	Name  string
	Logo  string
	Level uuid.UUID
}
