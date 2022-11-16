package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type EventMinimal interface {
	ID() uuid.UUID
	Title() string
	Organizer() uuid.UUID
	SubmissionDeadline() time.Time
}

type EventService interface {
	GetAll(ctx context.Context) ([]models.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Event, error)
	Create(ctx context.Context, event models.Event) (models.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetAllAsMinimal(ctx context.Context, id uuid.UUID) ([]EventMinimal, error)
	GetByIDAsMinimal(ctx context.Context, id uuid.UUID) (EventMinimal, error)
}
