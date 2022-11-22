package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type OrganizerService interface {
	GetAll(ctx context.Context) ([]models.Organizer, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Organizer, error)
	Create(ctx context.Context, name, logo string, level uuid.UUID) (models.Organizer, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, name, logo string, level uuid.UUID) (models.Organizer, error)

	GetAllLevels(ctx context.Context) ([]models.OrganizerLevel, error)
	CreateLevel(ctx context.Context, name string, code string) (models.OrganizerLevel, error)
	UpdateLevel(ctx context.Context, level models.OrganizerLevel) (models.OrganizerLevel, error)
}
