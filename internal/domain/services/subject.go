package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type SubjectService interface {
	GetAllExisting(ctx context.Context) ([]models.Subject, error)
	GetAllForEvent(ctx context.Context, eventId uuid.UUID) ([]models.Subject, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Subject, error)
	Create(ctx context.Context, eventId uuid.UUID, subject string) (models.Subject, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, subject models.Subject) error
}
