package storages

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type SubjectStorageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.Subject, error)
	GetByEvent(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
	GetAll(ctx context.Context) ([]models.Subject, error)
	Add(ctx context.Context, subject models.Subject) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, subject models.Subject) error
}
