package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type EventMinimal struct {
	ID                 uuid.UUID
	Title              string
	Organizer          uuid.UUID
	SubmissionDeadline time.Time
}

type EventCreateUpdateInfo struct {
	Title               string
	Organizer           uuid.UUID
	FoundingType        string
	FoundingRangeLow    int
	FoundingRangeHigh   int
	CoFoundingRangeLow  int
	CoFoundingRangeHigh int
	SubmissionDeadline  time.Time
	ConsiderationPeriod string
	RealisationPeriod   string
	Result              string
	Site                string
	Document            string
	InternalContacts    string
	TRL                 int
	Competitors         []uuid.UUID
	Subjects            []string
}

type EventService interface {
	GetAll(ctx context.Context) ([]models.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Event, error)
	Create(ctx context.Context, info EventCreateUpdateInfo) (models.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetAllAsMinimal(ctx context.Context) ([]EventMinimal, error)
	GetByIDAsMinimal(ctx context.Context, id uuid.UUID) (EventMinimal, error)
}
