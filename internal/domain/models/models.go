package models

import (
	"time"

	"github.com/google/uuid"
)

type OrganizerLevel struct {
	ID   uuid.UUID
	Name string
	Code string
}

type Organizer struct {
	ID    uuid.UUID
	Name  string
	Logo  string
	Level uuid.UUID
}

type RangeModel struct {
	ID   uuid.UUID
	Low  int
	High int
}

type FoundingRange = RangeModel
type CoFoundingRange = RangeModel

type Competitor struct {
	ID   uuid.UUID
	Name string
}

type Event struct {
	ID                  uuid.UUID
	Title               string
	Organizer           uuid.UUID
	FoundingType        string
	FoundingRange       uuid.UUID
	CoFoundingRange     uuid.UUID
	SubmissionDeadline  time.Time
	ConsiderationPeriod string
	RealisationPeriod   string
	Result              string
	Site                string
	Document            string
	InternalContacts    string
	TRL                 int
	Competitors         []uuid.UUID
}

type Subject struct {
	ID      uuid.UUID
	Name    string
	EventID uuid.UUID
}

type StoredImage struct {
	Link  string
	Value []byte
}

type User struct {
	ID       uuid.UUID
	Name     string
	Password string
}
