package models

import "github.com/google/uuid"

type OrganizerLevel struct {
	ID   uuid.UUID
	Name string
	Code string
}
