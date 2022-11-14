package storages

import (
	"context"

	"github.com/google/uuid"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

type RangeStorageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.FoundingRange, error)
	GetMaximumRange(ctx context.Context) (models.FoundingRange, error)
	Create(ctx context.Context, foundingRange models.FoundingRange) (models.FoundingRange, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
