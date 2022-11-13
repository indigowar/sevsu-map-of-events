package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type OrganizerStorageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.Organizer, error)
	GetAll(ctx context.Context) ([]models.Organizer, error)
	Create(ctx context.Context, organizer models.Organizer) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, organizer models.Organizer) error

	GetLevels(ctx context.Context) ([]models.OrganizerLevel, error)
	AddLevel(ctx context.Context, level models.OrganizerLevel) error
	DeleteLevel(ctx context.Context, id uuid.UUID, replacement uuid.UUID) error
}
