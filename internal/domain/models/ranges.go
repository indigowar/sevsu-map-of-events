package models

import "github.com/google/uuid"

type RangeModel struct {
	ID   uuid.UUID
	Low  int
	High int
}

type FoundingRange = RangeModel
type CoFoundingRange = RangeModel
